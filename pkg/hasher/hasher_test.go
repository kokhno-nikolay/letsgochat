package hasher_test

import (
	"fmt"
	"testing"

	"github.com/kokhno-nikolay/letsgochat/pkg/hasher"
)

func TestHashPassword(t *testing.T) {
	for i, tt := range []struct {
		pass   string
		secret string
		hash   string
	}{
		{
			"5bf07105-6aac-4220-a5bc-c90590b38a20",
			"47H37P6dRfN66DLy5rCA3sP37xdzdXkh",
			"5853fe150561d796d8e10e301db899af0000000000000000000000000000000000000000",
		},
		{
			"c1f8070d-0c49-434f-b523-cf8c9c02eef6",
			"di542eX9LzYA38xaH59MhT7Cr4v9cBsP",
			"c0c76d33ca78cdad7be0df5d889181ab0000000000000000000000000000000000000000",
		},
		{
			"7dd1c7a5-1499-43a4-9a0f-b100f8602581",
			"thisis32bitlongpassphraseimusing",
			"f4e6aa51034fa0df746598a470e571f90000000000000000000000000000000000000000",
		},
		{
			"48946012-29bf-4e06-8cf6-44be1cdc41a5",
			"B88VhC7CR5bS33m4E8Pr9ZM4zkj9ykne",
			"100cbd3b57c5673167389c8cbca2fd850000000000000000000000000000000000000000",
		},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			hasher := hasher.NewHasher(tt.secret)
			hash, err := hasher.HashPassword(tt.pass)
			if err != nil {
				t.Fatalf(err.Error())
			}

			if hash != tt.hash {
				t.Errorf("want %v; got %v", tt.hash, hash)
			}
		})
	}
}

func TestCheckHashPassword(t *testing.T) {
	for i, tt := range []struct {
		pass   string
		secret string
		hash   string
		res    bool
	}{
		{
			"5bf07105-6aac-4220-a5bc-c90590b38a20",
			"47H37P6dRfN66DLy5rCA3sP37xdzdXkh",
			"5853fe150561d796d8e10e301db899af0000000000000000000000000000000000000000",
			true,
		},
		{
			"5bf07105-6aac-4220-a5bc-c90590b38a20",
			"47H37P6dRfN66DLy5rCA3sP37xdzdXkh",
			"bad-hash",
			false,
		},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			h := hasher.NewHasher(tt.secret)

			res := h.CheckHashPassword(tt.pass, tt.hash)
			if res != tt.res {
				t.Errorf("want %v; got %v", tt.res, res)
			}
		})
	}
}

func BenchmarkHashPassword(b *testing.B) {
	pass := "5bf07105-6aac-4220-a5bc-c90590b38a20"
	secret := "47H37P6dRfN66DLy5rCA3sP37xdzdXkh"
	h := hasher.NewHasher(secret)

	for n := 0; n < b.N; n++ {
		h.HashPassword(pass)
	}
}

func BenchmarkCheckHashPassword(b *testing.B) {
	pass := "5bf07105-6aac-4220-a5bc-c90590b38a20"
	secret := "47H37P6dRfN66DLy5rCA3sP37xdzdXkh"
	hash := "5853fe150561d796d8e10e301db899af0000000000000000000000000000000000000000"
	h := hasher.NewHasher(secret)

	for n := 0; n < b.N; n++ {
		h.CheckHashPassword(pass, hash)
	}
}
