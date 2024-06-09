package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}
func TestAdd(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	parcels, err := store.GetParcelByID(id)
	parcel.Number = parcels.Number
	assert.NoError(t, err)
	assert.Equal(t, parcel, parcel)

	err = store.Delete(id)
	require.NoError(t, err)
	_, err = store.GetParcelByID(id)
	if err != nil {
		require.NoError(t, err)
	}

}

func TestGet(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	parcels, err := store.GetByClient(parcel.Client)
	require.NoError(t, err)
	fmt.Println(parcels)
}

func TestDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	_, err = store.GetByClient(parcel.Client)
	require.NoError(t, err)
	err = store.Delete(parcel.Client)
	require.NoError(t, err)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	par := getTestParcel()
	id, err := store.Add(par)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)
	line, err := store.GetParcelByID(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, line.Address)

}
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = int(id)
		parcelMap[int(id)] = parcels[i]
	}
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Equal(t, len(parcels), len(storedParcels))
	require.NotEmpty(t, storedParcels)
	for _, parcels := range storedParcels {
		mapParcel := parcelMap[parcels.Number]
		require.Equal(t, parcels, mapParcel)
	}
}
