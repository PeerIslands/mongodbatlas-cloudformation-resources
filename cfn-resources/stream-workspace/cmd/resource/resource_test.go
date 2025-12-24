// Copyright 2024 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"testing"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"
)

func TestGetTierValue(t *testing.T) {
	testCases := []struct {
		name     string
		tier     string
		expected int
	}{
		{
			name:     "SP2 tier",
			tier:     "SP2",
			expected: 2,
		},
		{
			name:     "SP5 tier",
			tier:     "SP5",
			expected: 5,
		},
		{
			name:     "SP10 tier",
			tier:     "SP10",
			expected: 10,
		},
		{
			name:     "SP30 tier",
			tier:     "SP30",
			expected: 30,
		},
		{
			name:     "SP50 tier",
			tier:     "SP50",
			expected: 50,
		},
		{
			name:     "Invalid tier",
			tier:     "INVALID",
			expected: 0,
		},
		{
			name:     "Empty tier",
			tier:     "",
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getTierValue(tc.tier)
			assert.Equal(t, tc.expected, result, "tier value should match expected")
		})
	}
}

func TestValidateTierComparison(t *testing.T) {
	testCases := []struct {
		name           string
		streamConfig   *StreamConfig
		expectedError  bool
		expectedStatus handler.Status
		expectedMsg    string
	}{
		{
			name: "Valid: MaxTierSize equals Tier",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP30"),
				MaxTierSize: stringPtr("SP30"),
			},
			expectedError: false,
		},
		{
			name: "Valid: MaxTierSize greater than Tier (SP30 < SP50)",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP30"),
				MaxTierSize: stringPtr("SP50"),
			},
			expectedError: false,
		},
		{
			name: "Valid: MaxTierSize greater than Tier (SP2 < SP10)",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP2"),
				MaxTierSize: stringPtr("SP10"),
			},
			expectedError: false,
		},
		{
			name: "Valid: MaxTierSize greater than Tier (SP5 < SP30)",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP5"),
				MaxTierSize: stringPtr("SP30"),
			},
			expectedError: false,
		},
		{
			name: "Valid: MaxTierSize greater than Tier (SP10 < SP50)",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP10"),
				MaxTierSize: stringPtr("SP50"),
			},
			expectedError: false,
		},
		{
			name: "Invalid: MaxTierSize less than Tier (SP50 < SP30)",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP50"),
				MaxTierSize: stringPtr("SP30"),
			},
			expectedError:  true,
			expectedStatus: handler.Failed,
			expectedMsg:    "MaxTierSize (SP30) must not be less than Tier (SP50)",
		},
		{
			name: "Invalid: MaxTierSize less than Tier (SP30 < SP10)",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP30"),
				MaxTierSize: stringPtr("SP10"),
			},
			expectedError:  true,
			expectedStatus: handler.Failed,
			expectedMsg:    "MaxTierSize (SP10) must not be less than Tier (SP30)",
		},
		{
			name: "Invalid: MaxTierSize less than Tier (SP10 < SP5)",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP10"),
				MaxTierSize: stringPtr("SP5"),
			},
			expectedError:  true,
			expectedStatus: handler.Failed,
			expectedMsg:    "MaxTierSize (SP5) must not be less than Tier (SP10)",
		},
		{
			name: "Invalid: MaxTierSize less than Tier (SP5 < SP2)",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP5"),
				MaxTierSize: stringPtr("SP2"),
			},
			expectedError:  true,
			expectedStatus: handler.Failed,
			expectedMsg:    "MaxTierSize (SP2) must not be less than Tier (SP5)",
		},
		{
			name: "Invalid: MaxTierSize less than Tier (SP50 < SP2)",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP50"),
				MaxTierSize: stringPtr("SP2"),
			},
			expectedError:  true,
			expectedStatus: handler.Failed,
			expectedMsg:    "MaxTierSize (SP2) must not be less than Tier (SP50)",
		},
		{
			name: "Valid: StreamConfig is nil",
			streamConfig: nil,
			expectedError: false,
		},
		{
			name: "Valid: Tier is nil",
			streamConfig: &StreamConfig{
				Tier:        nil,
				MaxTierSize: stringPtr("SP30"),
			},
			expectedError: false,
		},
		{
			name: "Valid: MaxTierSize is nil",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP30"),
				MaxTierSize: nil,
			},
			expectedError: false,
		},
		{
			name: "Valid: Both Tier and MaxTierSize are nil",
			streamConfig: &StreamConfig{
				Tier:        nil,
				MaxTierSize: nil,
			},
			expectedError: false,
		},
		{
			name: "Valid: All enum values in order",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP2"),
				MaxTierSize: stringPtr("SP5"),
			},
			expectedError: false,
		},
		{
			name: "Valid: All enum values in order - next pair",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP5"),
				MaxTierSize: stringPtr("SP10"),
			},
			expectedError: false,
		},
		{
			name: "Valid: All enum values in order - next pair",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP10"),
				MaxTierSize: stringPtr("SP30"),
			},
			expectedError: false,
		},
		{
			name: "Valid: All enum values in order - last pair",
			streamConfig: &StreamConfig{
				Tier:        stringPtr("SP30"),
				MaxTierSize: stringPtr("SP50"),
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validateTierComparison(tc.streamConfig)

			if tc.expectedError {
				assert.NotNil(t, result, "expected validation error but got nil")
				if result != nil {
					assert.Equal(t, tc.expectedStatus, result.OperationStatus, "operation status should match")
					assert.Contains(t, result.Message, tc.expectedMsg, "error message should contain expected text")
					assert.Equal(t, string(types.HandlerErrorCodeInvalidRequest), result.HandlerErrorCode, "error code should be InvalidRequest")
				}
			} else {
				assert.Nil(t, result, "expected no validation error but got error: %v", result)
			}
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

