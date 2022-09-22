package orm

import "strconv"

// ====== 原生类型数据 转 字节数组 ======

func boolToBytes(val bool) []byte {
	return strconv.AppendBool(nil, val)
}

func intToBytes[T FieldTypeInt](val T) []byte {
	return strconv.AppendInt(nil, int64(val), 10)
}

func uintToBytes[T FieldTypeUint](val T) []byte {
	return strconv.AppendUint(nil, uint64(val), 10)
}

func floatToBytes[T FieldTypeFloat](val T, bitSize int) []byte {
	// FIXME:  fmt=g 会导致过大数值变成科学计数
	// return strconv.AppendFloat(nil, float64(val), 'g', -1, bitSize)
	return strconv.AppendFloat(nil, float64(val), 'f', -1, bitSize)
}

// ====== 字节数组 转 原生类型数据 ======

func bytesToBool(body []byte) (val bool, err error) {
	return strconv.ParseBool(string(body))
}

func bytesToInt[T FieldTypeInt](body []byte) (val T, err error) {
	var vv int64
	vv, err = strconv.ParseInt(string(body), 10, 0)
	if err != nil {
		return
	}

	val = T(vv)
	return
}

func bytesToUint[T FieldTypeUint](body []byte) (val T, err error) {
	var vv uint64
	vv, err = strconv.ParseUint(string(body), 10, 0)
	if err != nil {
		return
	}

	val = T(vv)
	return
}

func bytesToFloat[T FieldTypeFloat](body []byte, bitSize int) (val T, err error) {
	var vv float64
	vv, err = strconv.ParseFloat(string(body), bitSize)
	if err != nil {
		return
	}

	val = T(vv)
	return
}
