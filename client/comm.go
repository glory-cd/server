/**
* @Author: xhzhang
* @Date: 2019/7/26 13:37
 */
package client

import (
	"errors"
	"strconv"
	"strings"
)

///*
//	判断服务ID列表是否是某个Group下的服务
//    返回值：bool : true 或者 false
//           []string : 不属于的服务ID列表
//*/
//func ServiceIsGroupSubSet(serviceIDs []string, groupServices map[string]int) (bool, []string) {
//	notbelong := []string{}
//	for _, sid := range serviceIDs {
//		if _, ok := groupServices[sid]; !ok {
//			notbelong = append(notbelong, sid)
//		}
//	}
//	if len(notbelong) > 0 {
//		return false, notbelong
//	}
//
//	return true, notbelong
//}
//
//func ReleaseCodeIsReleaseSubSet(releaseCodeIDs []int, releaseCodes map[int]string) (bool, []string) {
//	notbelong := []string{}
//	for _, rcid := range releaseCodeIDs {
//		if _, ok := releaseCodes[rcid]; !ok {
//			notbelong = append(notbelong,strconv.Itoa(rcid))
//		}
//	}
//	if len(notbelong) > 0 {
//		return false, notbelong
//	}
//
//	return true, notbelong
//}
//
func IsSubSet(IDs []interface{}, pMap map[interface{}]interface{}) (bool, []string) {
	notbelong := []string{}
	for _, id := range IDs {
		if _, ok := pMap[id]; !ok {
			switch id.(type) {
			case int:
				notbelong = append(notbelong, strconv.Itoa(id.(int)))
			case string:
				notbelong = append(notbelong, id.(string))
			}
		}
	}
	if len(notbelong) > 0 {
		return false, notbelong
	}

	return true, notbelong
}

/*
	改变service对象形式为map[string]int。key: 服务ID，val: 固定1
*/
func GenerateServiceMap(services []Service) map[interface{}]interface{} {
	serviceParentSet := make(map[interface{}]interface{})
	for _, s := range services {
		serviceParentSet[s.ID] = 1
	}
	return serviceParentSet
}

func String2Interface(strs []string) (iStrs []interface{}) {
	for _, n := range strs {
		iStrs = append(iStrs, n)
	}
	return iStrs
}

func Int2Interface(ints []int32) (iInts []interface{}) {
	for _, n := range ints {
		iInts = append(iInts, n)
	}
	return iInts
}

/*
  判断服务是否属于指定group
*/
func CheckServiceOwnGroup(IDs []string, groupServices []Service) error {
	serviceParentSet := GenerateServiceMap(groupServices)
	isServiceLegal, serviceNotBelong := IsSubSet(String2Interface(IDs), serviceParentSet)
	if !isServiceLegal {
		return errors.New("service ID not match with group:" + strings.Join(serviceNotBelong, ";"))
	}
	return nil
}

func CheckRCIDOwnRelease(IDs []int32, rcMap map[string]int32) error {
	nMap := make(map[interface{}]interface{})
	for k, v := range rcMap {
		nMap[v] = k
	}
	isReleaseCodeLegal, releaseCodeNotBelong := IsSubSet(Int2Interface(IDs), nMap)
	if !isReleaseCodeLegal {
		return errors.New("发布代码ID不匹配:" + strings.Join(releaseCodeNotBelong, ";"))
	}
	return nil
}

func CheckRCNameOwnRelease(names []string, rcMap map[string]int32) error {
	reverseMap := make(map[interface{}]interface{})
	for k, v := range rcMap {
		reverseMap[v] = k
	}

	isSub, notBelong := IsSubSet(String2Interface(names), reverseMap)
	if !isSub {
		return errors.New("代码名称与本次发布不匹配:" + strings.Join(notBelong, ";"))
	}
	return nil
}
