package provider

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"terraform-provider-xrcm/internal/xrcm_pf"

	"github.com/fujiwara/tfstate-lookup/tfstate"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func getContent(body []byte) (c map[string]interface{}, err error) {
	var data = make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var content map[string]interface{}
	content, contentOk := data["content"].(map[string]interface{})
	if !contentOk {
		return nil, errors.New("getResourceIdNContent: No Content data:")
	} else {
		return content, nil
	}
}

func getResourceIdNContent(body []byte) (r map[string]interface{}, c map[string]interface{}, err error) {
	var data1 = make(map[string]interface{})
	err = json.Unmarshal(body, &data1)
	if err != nil {
		return nil, nil, err
	}

	data, ok := data1["data"].(map[string]interface{})
	if !ok {
		return nil, nil, err
	}
	var resourceId map[string]interface{}
	resourceId, resourceIdOk := data["resourceId"].(map[string]interface{})

	var content map[string]interface{}
	content, contentOk := data["content"].(map[string]interface{})

	if !contentOk && !resourceIdOk {
		return nil, nil, errors.New("getResourceIdNContent: No ResourceID and Content data:")
	} else if !contentOk {
		return resourceId, nil, errors.New("getResourceIdNContent: No Content data")
	} else if !resourceIdOk {
		return nil, content, errors.New("getResourceIdNContent: No ResourceID data")
	} else {
		return resourceId, content, nil
	}
}

// func GetData(plan *interface{}, body []byte) (c map[string]interface{}, err error) {
func SetResourceId(deviceName string, Id *types.String, body []byte) (c map[string]interface{}, err error) {

	var resourceId map[string]interface{}
	var content map[string]interface{}
	resourceId, content, _ = getResourceIdNContent(body)

	if resourceId == nil && content == nil {
		return nil, errors.New("SetResourceId: No ResourceID and/or Content data:" + string(body))
	}

	if Id != nil && len(Id.ValueString()) <= 0 {
		if content["href"] == nil {
			*Id = types.StringValue(deviceName + resourceId["href"].(string))
		} else {
			*Id = types.StringValue(deviceName + content["href"].(string))
		}
	}
	return content, nil
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[pos:]
}

var mytfstate *tfstate.TFState = nil
var tfstatefile string = ""

func GetTFState(ctx context.Context, file string) (*tfstate.TFState, error) {
	if mytfstate == nil || tfstatefile != file {
		tfstatefile = file
		mytfstate, _ = tfstate.ReadFile(ctx, file)
	}
	return mytfstate, nil
}

func LookupTFState(key string) (interface{}, error) {
	if mytfstate != nil {
		value, _ := mytfstate.Lookup(key)
		if value != nil {
			return value.Value, nil
		}
	}
	return nil, nil
}

func GetAndLookupTFState(ctx context.Context, file string, key string) (interface{}, error) {
	state, _ := GetTFState(ctx, file)
	if state != nil {
		value, _ := state.Lookup(key)
		if value != nil {
			return value.Value, nil
		}
	}
	return nil, nil
}

func GetResource(ctx context.Context, client *xrcm_pf.Client, deviceName string, query string) (map[string]interface{}, string, error) {

	tflog.Debug(ctx, "GetResource: ", map[string]interface{}{"queryData": query})

	body, deviceId, err := client.ExecuteDeviceHttpCommand(deviceName, "GET", query, nil)

	if err != nil {
		return nil, "", errors.New("GetResource: Could not query, unexpected error:" + err.Error())
	}

	tflog.Debug(ctx, "GetResource: Query SUCCESS", map[string]interface{}{"body": string(body)})

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, deviceId, errors.New("GetResource: Could not parse query body, unexpected error:" + err.Error())
	}
	tflog.Debug(ctx, "GetResource: get resource", map[string]interface{}{"data": data})
	return data, deviceId, nil
}
func GetResourcebyID(ctx context.Context, client *xrcm_pf.Client, deviceId string, query string) (map[string]interface{}, error) {

	tflog.Debug(ctx, "GetResource: ", map[string]interface{}{"queryData": query})

	body, err := client.ExecuteDeviceHttpCommandByID(deviceId, "GET", query, nil)

	if err != nil {
		return nil, errors.New("GetResource: Could not query, unexpected error:" + err.Error())
	}

	tflog.Debug(ctx, "GetResource: Query SUCCESS", map[string]interface{}{"body": string(body)})

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.New("GetResource: Could not parse query body, unexpected error:" + err.Error())
	}
	tflog.Debug(ctx, "GetResource: get resource", map[string]interface{}{"data": data})
	return data, nil
}

func Find(what string, where []types.String) (idx int) {
	for i, v := range where {
		if v.ValueString() == what {
			return i
		}
	}
	return -1
}

// Sets the bit at pos in the integer n.
func setBit(n int, pos uint) int {
	n |= (1 << pos)
	return n
}

// Clears the bit at pos in n.
func clearBit(n int, pos uint) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}

func hasBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func setBits(positions []int) int {
	n := 0
	for _, v := range positions {
		n = setBit(n, uint(v))
	}
	return n
}

func getBits(n int) []attr.Value {
	var bits []attr.Value
	for i := 0; i < 16; i++ {
		if hasBit(n, uint(i)) {
			bits = append(bits, types.Int64Value(int64(i)))
		}
	}
	return bits
}
