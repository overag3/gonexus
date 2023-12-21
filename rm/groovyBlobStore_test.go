package nexusrm

import (
	"context"
	"testing"
)

func TestCreateFileBlobStore(t *testing.T) {
	t.Skip("Needs new framework")
	rm, mock := repositoriesTestRM(t)
	defer mock.Close()

	err := CreateFileBlobStoreContext(context.Background(), rm, "testname", "testpath")
	if err != nil {
		t.Error(err)
	}

	// TODO: list blobstores
}

func TestCreateBlobStoreGroup(t *testing.T) {
	t.Skip("Needs new framework")
	rm, mock := repositoriesTestRM(t)
	defer mock.Close()

	CreateFileBlobStoreContext(context.Background(), rm, "f1", "pathf1")
	CreateFileBlobStoreContext(context.Background(), rm, "f2", "pathf2")
	CreateFileBlobStoreContext(context.Background(), rm, "f3", "pathf3")

	err := CreateBlobStoreGroupContext(context.Background(), rm, "grpname", []string{"f1", "f2", "f3"})
	if err != nil {
		t.Error(err)
	}
}

/*
func TestDeleteBlobStore(t *testing.T) {
	t.Skip("Needs new framework")
	rm := getTestRM(t)

	err := DeleteBlobStore(rm, "testname")
	if err != nil {
		t.Error(err)
	}

	// TODO: list blobstores
}
*/
