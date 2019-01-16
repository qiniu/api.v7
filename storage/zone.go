package storage

// Zone 是Region的别名
// 兼容保留
type Zone = Region

// GetZone 用来根据ak和bucket来获取空间相关的机房信息
// 新版本使用GetRegion, 这个函数用来保持兼容
func GetZone(ak, bucket string) (zone *Zone, err error) {
	return GetRegion(ak, bucket)
}

var (
	// 华东机房
	// 兼容保留
	ZoneHuadong, _ = GetDefaultRegion(RCodeHuadong)

	// 华北机房
	// 兼容保留
	ZoneHuabei, _ = GetDefaultRegion(RCodeHuabei)

	// 华南机房
	// 兼容保留
	ZoneHuanan, _ = GetDefaultRegion(RCodeHuanan)

	// 北美机房
	// 兼容保留
	ZoneBeimei, _ = GetDefaultRegion(RCodeBeimei)

	// 新加坡机房
	// 兼容保留
	ZoneXinjiapo, _ = GetDefaultRegion(RCodeAsia)

	// 兼容保留
	Zone_z0 = ZoneHuadong
	// 兼容保留
	Zone_z1 = ZoneHuabei
	// 兼容保留
	Zone_z2 = ZoneHuanan
	// 兼容保留
	Zone_na0 = ZoneBeimei
	// 兼容保留
	Zone_as0 = ZoneXinjiapo
)
