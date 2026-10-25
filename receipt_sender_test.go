package main

import "testing"

func TestHandoffMessage(t *testing.T) {
	tests := []struct {
		name  string
		order Order
		want  string
	}{
		{"seller asset included", Order{"ORD-7", "North Market", "+15550001", "oak desk"}, "Order ORD-7: North Market is ready for handoff. Seller asset: oak desk."},
		{"buyer update is specific", Order{"ORD-8", "Harbor Goods", "+15550002", "camera"}, "Order ORD-8: Harbor Goods is ready for handoff. Seller asset: camera."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandoffMessage(tt.order); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
