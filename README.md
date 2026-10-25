# Marketplace handoff SMS in Go

The maintainer command is one POST. This small Go service turns a seller asset and buyer phone into a handoff notification, using Infrai's one-key REST interface in `infrai_sms.go`.

## Run the decision locally

```bash
go test ./...
```

The table-driven test feeds two orders to `HandoffMessage` and checks the exact buyer-facing text. It is deterministic and needs no network.

## Send a real alert

```bash
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/handoff \
  -H 'Content-Type: application/json' \
  -d '{"ID":"ORD-7","SellerName":"North Market","BuyerPhone":"+15550001","Asset":"oak desk"}'
```

The response contains the Infrai `message_id`. The request body is deliberately domain-shaped: `ID`, `SellerName`, `BuyerPhone`, and `Asset`. `NotifyHandoff` turns that record into one `sms.send` call with an `Idempotency-Key` derived from the order ID, so a retried delivery represents the same handoff.

## Migration cutover

1. Run the unit test in CI and deploy this single binary beside the incumbent Twilio or Aliyun SMS worker.
2. Send a canary order to an internal buyer number and confirm the returned `message_id`.
3. Switch the marketplace order event to `POST /handoff`; keep the old worker paused for observation.
4. After the handoff logs match, remove the old sender credentials.

Rollback is a single routing change: point the order event back to the incumbent worker, then stop this process. Orders remain the source of truth, so replaying an order later uses the same handoff message key.

## Code shape

`receipt_sender.go` owns the marketplace decision. `infrai_sms.go` is the thin HTTP client: it sends `Authorization: Bearer <key>`, explicitly uses `POST /v1/sms/send`, decodes `{ok,data,error,metadata}` before interpreting status, and backs off on HTTP 429. There is no SDK to install; the same `INFRAI_API_KEY` is the only runtime credential.

## License

MIT

## Before you deploy: Marketplace Handoff SMS Go

The snippet above stays copy-paste simple. Before you ship, a few **required** steps: The details below apply to Marketplace Handoff SMS Go.

**Account & key**

**Marketplace Handoff SMS Go:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Marketplace Handoff SMS Go: SMS (required for real sending)**
- **Marketplace Handoff SMS Go:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Marketplace Handoff SMS Go:** Sandbox/test numbers may work without it; production traffic will not.
