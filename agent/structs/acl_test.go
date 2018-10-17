package structs

import (
	"strings"
	"testing"

	"github.com/hashicorp/consul/acl"

	"github.com/stretchr/testify/require"
)

func TestStructs_ACLToken_PolicyIDs(t *testing.T) {
	t.Parallel()

	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		token := &ACLToken{
			Policies: []ACLTokenPolicyLink{
				ACLTokenPolicyLink{
					ID: "one",
				},
				ACLTokenPolicyLink{
					ID: "two",
				},
				ACLTokenPolicyLink{
					ID: "three",
				},
			},
		}

		policyIDs := token.PolicyIDs()
		require.Len(t, policyIDs, 3)
		require.Equal(t, "one", policyIDs[0])
		require.Equal(t, "two", policyIDs[1])
		require.Equal(t, "three", policyIDs[2])
	})

	t.Run("Legacy Management", func(t *testing.T) {
		t.Parallel()

		a := &ACL{
			ID:   "root",
			Type: ACLTokenTypeManagement,
			Name: "management",
		}

		token := a.Convert()

		policyIDs := token.PolicyIDs()
		require.Len(t, policyIDs, 1)
		require.Equal(t, ACLPolicyGlobalManagementID, policyIDs[0])
	})

	t.Run("No Policies", func(t *testing.T) {
		t.Parallel()

		token := &ACLToken{}

		policyIDs := token.PolicyIDs()
		require.Len(t, policyIDs, 0)
	})
}

func TestStructs_ACLToken_EmbeddedPolicy(t *testing.T) {
	t.Parallel()

	t.Run("No Rules", func(t *testing.T) {
		t.Parallel()

		token := &ACLToken{}
		require.Nil(t, token.EmbeddedPolicy())
	})

	t.Run("Legacy Client", func(t *testing.T) {
		t.Parallel()

		// None of the other fields should be considered
		token := &ACLToken{
			Type:  ACLTokenTypeClient,
			Rules: `acl = "read"`,
		}

		policy := token.EmbeddedPolicy()
		require.NotNil(t, policy)
		require.NotEqual(t, "", policy.ID)
		require.True(t, strings.HasPrefix(policy.Name, "legacy-policy-"))
		require.Equal(t, token.Rules, policy.Rules)
		require.Equal(t, policy.Syntax, acl.SyntaxLegacy)
		require.NotNil(t, policy.Hash)
		require.NotEqual(t, []byte{}, policy.Hash)
	})

	t.Run("Same Policy for Tokens with same Rules", func(t *testing.T) {
		t.Parallel()

		token1 := &ACLToken{
			AccessorID:  "f55b260c-5e05-418e-ab19-d421d1ab4b52",
			SecretID:    "b2165bac-7006-459b-8a72-7f549f0f06d6",
			Description: "token 1",
			Type:        ACLTokenTypeClient,
			Rules:       `acl = "read"`,
		}

		token2 := &ACLToken{
			AccessorID:  "09d1c059-961a-46bd-a2e4-76adebe35fa5",
			SecretID:    "65e98e67-9b29-470c-8ffa-7c5a23cc67c8",
			Description: "token 2",
			Type:        ACLTokenTypeClient,
			Rules:       `acl = "read"`,
		}

		policy1 := token1.EmbeddedPolicy()
		policy2 := token2.EmbeddedPolicy()
		require.Equal(t, policy1, policy2)
	})
}

func TestStructs_ACLToken_SetHash(t *testing.T) {
	t.Parallel()

	token := ACLToken{
		AccessorID:  "09d1c059-961a-46bd-a2e4-76adebe35fa5",
		SecretID:    "65e98e67-9b29-470c-8ffa-7c5a23cc67c8",
		Description: "test",
		Policies: []ACLTokenPolicyLink{
			ACLTokenPolicyLink{
				ID: "one",
			},
			ACLTokenPolicyLink{
				ID: "two",
			},
			ACLTokenPolicyLink{
				ID: "three",
			},
		},
	}

	t.Run("Nil Hash - Generate", func(t *testing.T) {
		require.Nil(t, token.Hash)
		h := token.SetHash(false)
		require.NotNil(t, h)
		require.NotEqual(t, []byte{}, h)
		require.Equal(t, h, token.Hash)
	})

	t.Run("Hash Set - Dont Generate", func(t *testing.T) {
		original := token.Hash
		h := token.SetHash(false)
		require.Equal(t, original, h)

		token.Description = "changed"
		h = token.SetHash(false)
		require.Equal(t, original, h)
	})

	t.Run("Hash Set - Generate", func(t *testing.T) {
		original := token.Hash
		h := token.SetHash(true)
		require.NotEqual(t, original, h)
	})
}

func TestStructs_ACLToken_EstimateSize(t *testing.T) {
	t.Parallel()

	// estimated size here should
	token := ACLToken{
		AccessorID:  "09d1c059-961a-46bd-a2e4-76adebe35fa5",
		SecretID:    "65e98e67-9b29-470c-8ffa-7c5a23cc67c8",
		Description: "test",
		Policies: []ACLTokenPolicyLink{
			ACLTokenPolicyLink{
				ID: "one",
			},
			ACLTokenPolicyLink{
				ID: "two",
			},
			ACLTokenPolicyLink{
				ID: "three",
			},
		},
	}

	// this test is very contrived. Basically just tests that the
	// math is okay and returns the value.
	require.Equal(t, 120, token.EstimateSize())
}

func TestStructs_ACLToken_Stub(t *testing.T) {
	t.Parallel()

	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		token := ACLToken{
			AccessorID:  "09d1c059-961a-46bd-a2e4-76adebe35fa5",
			SecretID:    "65e98e67-9b29-470c-8ffa-7c5a23cc67c8",
			Description: "test",
			Policies: []ACLTokenPolicyLink{
				ACLTokenPolicyLink{
					ID: "one",
				},
				ACLTokenPolicyLink{
					ID: "two",
				},
				ACLTokenPolicyLink{
					ID: "three",
				},
			},
		}

		stub := token.Stub()

		require.Equal(t, token.AccessorID, stub.AccessorID)
		require.Equal(t, token.Description, stub.Description)
		require.Equal(t, token.Policies, stub.Policies)
		require.Equal(t, token.Local, stub.Local)
		require.Equal(t, token.CreateTime, stub.CreateTime)
		require.Equal(t, token.Hash, stub.Hash)
		require.Equal(t, token.CreateIndex, stub.CreateIndex)
		require.Equal(t, token.ModifyIndex, stub.ModifyIndex)
		require.False(t, stub.Legacy)
	})

	t.Run("Legacy", func(t *testing.T) {
		t.Parallel()
		token := ACLToken{
			AccessorID:  "09d1c059-961a-46bd-a2e4-76adebe35fa5",
			SecretID:    "65e98e67-9b29-470c-8ffa-7c5a23cc67c8",
			Description: "test",
			Type:        ACLTokenTypeClient,
			Rules:       `key "" { policy = "read" }`,
		}

		stub := token.Stub()
		require.Equal(t, token.AccessorID, stub.AccessorID)
		require.Equal(t, token.Description, stub.Description)
		require.Equal(t, token.Policies, stub.Policies)
		require.Equal(t, token.Local, stub.Local)
		require.Equal(t, token.CreateTime, stub.CreateTime)
		require.Equal(t, token.Hash, stub.Hash)
		require.Equal(t, token.CreateIndex, stub.CreateIndex)
		require.Equal(t, token.ModifyIndex, stub.ModifyIndex)
		require.True(t, stub.Legacy)
	})
}

func TestStructs_ACLTokens_Sort(t *testing.T) {
	t.Parallel()

	tokens := ACLTokens{
		&ACLToken{
			AccessorID: "9db509a9-c809-48c1-895d-99f845b7a9d5",
		},
		&ACLToken{
			AccessorID: "6bd01084-1695-43b8-898d-b2dd7874754d",
		},
		&ACLToken{
			AccessorID: "614a4cef-9149-4271-b878-7edb1ad661f8",
		},
		&ACLToken{
			AccessorID: "c9dd9980-8d54-472f-9e5e-74c02143e1f4",
		},
	}

	tokens.Sort()
	require.Equal(t, tokens[0].AccessorID, "614a4cef-9149-4271-b878-7edb1ad661f8")
	require.Equal(t, tokens[1].AccessorID, "6bd01084-1695-43b8-898d-b2dd7874754d")
	require.Equal(t, tokens[2].AccessorID, "9db509a9-c809-48c1-895d-99f845b7a9d5")
	require.Equal(t, tokens[3].AccessorID, "c9dd9980-8d54-472f-9e5e-74c02143e1f4")
}

func TestStructs_ACLTokenListStubs_Sort(t *testing.T) {
	t.Parallel()

	tokens := ACLTokenListStubs{
		&ACLTokenListStub{
			AccessorID: "9db509a9-c809-48c1-895d-99f845b7a9d5",
		},
		&ACLTokenListStub{
			AccessorID: "6bd01084-1695-43b8-898d-b2dd7874754d",
		},
		&ACLTokenListStub{
			AccessorID: "614a4cef-9149-4271-b878-7edb1ad661f8",
		},
		&ACLTokenListStub{
			AccessorID: "c9dd9980-8d54-472f-9e5e-74c02143e1f4",
		},
	}

	tokens.Sort()
	require.Equal(t, tokens[0].AccessorID, "614a4cef-9149-4271-b878-7edb1ad661f8")
	require.Equal(t, tokens[1].AccessorID, "6bd01084-1695-43b8-898d-b2dd7874754d")
	require.Equal(t, tokens[2].AccessorID, "9db509a9-c809-48c1-895d-99f845b7a9d5")
	require.Equal(t, tokens[3].AccessorID, "c9dd9980-8d54-472f-9e5e-74c02143e1f4")
}

func TestStructs_ACLPolicy_Stub(t *testing.T) {
	t.Parallel()

	policy := &ACLPolicy{
		ID:          "09d1c059-961a-46bd-a2e4-76adebe35fa5",
		Name:        "test",
		Description: "test",
		Rules:       `acl = "read"`,
	}

	stub := policy.Stub()

	require.Equal(t, policy.ID, stub.ID)
	require.Equal(t, policy.Name, stub.Name)
	require.Equal(t, policy.Description, stub.Description)
	require.Equal(t, policy.Datacenters, stub.Datacenters)
	require.Equal(t, policy.Hash, stub.Hash)
	require.Equal(t, policy.CreateIndex, stub.CreateIndex)
	require.Equal(t, policy.ModifyIndex, stub.ModifyIndex)
}

func TestStructs_ACLPolicy_SetHash(t *testing.T) {
	t.Parallel()

	policy := &ACLPolicy{
		ID:          "09d1c059-961a-46bd-a2e4-76adebe35fa5",
		Name:        "test",
		Description: "test",
		Rules:       `acl = "read"`,
	}

	t.Run("Nil Hash - Generate", func(t *testing.T) {
		require.Nil(t, policy.Hash)
		h := policy.SetHash(false)
		require.NotNil(t, h)
		require.NotEqual(t, []byte{}, h)
		require.Equal(t, h, policy.Hash)
	})

	t.Run("Hash Set - Dont Generate", func(t *testing.T) {
		original := policy.Hash
		h := policy.SetHash(false)
		require.Equal(t, original, h)

		policy.Description = "changed"
		h = policy.SetHash(false)
		require.Equal(t, original, h)
	})

	t.Run("Hash Set - Generate", func(t *testing.T) {
		original := policy.Hash
		h := policy.SetHash(true)
		require.NotEqual(t, original, h)
	})
}

func TestStructs_ACLPolicy_EstimateSize(t *testing.T) {
	t.Parallel()

	policy := ACLPolicy{
		ID:          "09d1c059-961a-46bd-a2e4-76adebe35fa5",
		Name:        "test",
		Description: "test",
		Rules:       `acl = "read"`,
	}

	// this test is very contrived. Basically just tests that the
	// math is okay and returns the value.
	require.Equal(t, 84, policy.EstimateSize())
	policy.Datacenters = []string{"dc1", "dc2"}
	require.Equal(t, 90, policy.EstimateSize())
}

func TestStructs_ACLPolicies_Sort(t *testing.T) {
	t.Parallel()

	policies := ACLPolicies{
		&ACLPolicy{
			ID: "9db509a9-c809-48c1-895d-99f845b7a9d5",
		},
		&ACLPolicy{
			ID: "6bd01084-1695-43b8-898d-b2dd7874754d",
		},
		&ACLPolicy{
			ID: "614a4cef-9149-4271-b878-7edb1ad661f8",
		},
		&ACLPolicy{
			ID: "c9dd9980-8d54-472f-9e5e-74c02143e1f4",
		},
	}

	policies.Sort()
	require.Equal(t, policies[0].ID, "614a4cef-9149-4271-b878-7edb1ad661f8")
	require.Equal(t, policies[1].ID, "6bd01084-1695-43b8-898d-b2dd7874754d")
	require.Equal(t, policies[2].ID, "9db509a9-c809-48c1-895d-99f845b7a9d5")
	require.Equal(t, policies[3].ID, "c9dd9980-8d54-472f-9e5e-74c02143e1f4")
}

func TestStructs_ACLPolicyListStubs_Sort(t *testing.T) {
	t.Parallel()

	policies := ACLPolicyListStubs{
		&ACLPolicyListStub{
			ID: "9db509a9-c809-48c1-895d-99f845b7a9d5",
		},
		&ACLPolicyListStub{
			ID: "6bd01084-1695-43b8-898d-b2dd7874754d",
		},
		&ACLPolicyListStub{
			ID: "614a4cef-9149-4271-b878-7edb1ad661f8",
		},
		&ACLPolicyListStub{
			ID: "c9dd9980-8d54-472f-9e5e-74c02143e1f4",
		},
	}

	policies.Sort()
	require.Equal(t, policies[0].ID, "614a4cef-9149-4271-b878-7edb1ad661f8")
	require.Equal(t, policies[1].ID, "6bd01084-1695-43b8-898d-b2dd7874754d")
	require.Equal(t, policies[2].ID, "9db509a9-c809-48c1-895d-99f845b7a9d5")
	require.Equal(t, policies[3].ID, "c9dd9980-8d54-472f-9e5e-74c02143e1f4")
}
