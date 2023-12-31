package pyrdp

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func ParseClientInfo(d *decode.D, length int64) {
	d.FieldStruct("client_info", func(d *decode.D) {
		pos := d.Pos()
		var (
			is_unicode bool
			has_null   bool
			null_n     uint64 = 0
			unicode_n  uint64 = 0
		)
		code_page := d.FieldU32("code_page")
		flags := d.U32()
		d.SeekRel(-4 * 8)
		d.FieldStruct("flags", decodeFlagsFn)

		is_unicode = ((flags & INFO_UNICODE) != 0)
		has_null = (code_page == 1252 || is_unicode)

		if has_null {
			null_n = 1
		}
		if is_unicode {
			unicode_n = 2
		}

		domain_length := int(d.FieldU16("domain_length") + null_n*unicode_n)
		username_length := int(d.FieldU16("username_length") + null_n*unicode_n)
		password_length := int(d.FieldU16("password_length") + null_n*unicode_n)
		alternate_shell_length := int(d.FieldU16("alternate_shell_length") + null_n*unicode_n)
		working_dir_length := int(d.FieldU16("working_dir_length") + null_n*unicode_n)

		d.FieldStrFn("domain", toTextUTF16Fn(domain_length))
		d.FieldStrFn("username", toTextUTF16Fn(username_length))
		d.FieldStrFn("password", toTextUTF16Fn(password_length))
		d.FieldStrFn("alternate_shell", toTextUTF16Fn(alternate_shell_length))
		d.FieldStrFn("working_dir", toTextUTF16Fn(working_dir_length))

		extra_length := length - ((d.Pos() - pos) / 8)
		if extra_length > 0 {
			d.FieldStruct("extra_info", func(d *decode.D) {
				d.FieldU16("address_family", scalar.UintHex)
				address_length := int(d.FieldU16("address_length"))
				d.FieldStrFn("address", toTextUTF16Fn(address_length))
				dir_length := int(d.FieldU16("dir_length"))
				d.FieldStrFn("dir", toTextUTF16Fn(dir_length))
			})

			// XXX: there's more extra info but here's everything we need from the
			//			client (other than UTC info)
		}
	})
}

const (
	// flags
	INFO_MOUSE                  = 0x00000001
	INFO_DISABLECTRLALTDEL      = 0x00000002
	INFO_AUTOLOGON              = 0x00000008
	INFO_UNICODE                = 0x00000010
	INFO_MAXIMIZESHELL          = 0x00000020
	INFO_LOGONNOTIFY            = 0x00000040
	INFO_COMPRESSION            = 0x00000080
	INFO_ENABLEWINDOWSKEY       = 0x00000100
	INFO_REMOTECONSOLEAUDIO     = 0x00002000
	INFO_FORCE_ENCRYPTED_CS_PDU = 0x00004000
	INFO_RAIL                   = 0x00008000
	INFO_LOGONERRORS            = 0x00010000
	INFO_MOUSE_HAS_WHEEL        = 0x00020000
	INFO_PASSWORD_IS_SC_PIN     = 0x00040000
	INFO_NOAUDIOPLAYBACK        = 0x00080000
	INFO_USING_SAVED_CREDS      = 0x00100000
	INFO_AUDIOCAPTURE           = 0x00200000
	INFO_VIDEO_DISABLE          = 0x00400000
	INFO_RESERVED1              = 0x00800000
	INFO_RESERVED2              = 0x01000000
	INFO_HIDEF_RAIL_SUPPORTED   = 0x02000000
)

func decodeFlagsFn(d *decode.D) {
	d.FieldBool("mouse")
	d.FieldBool("disabledctrlaltdel")
	d.SeekRel(1)
	d.FieldBool("autologon")
	d.FieldBool("unicode")
	d.FieldBool("maximizeshell")
	d.FieldBool("logonnotify")
	d.FieldBool("compression")
	d.FieldBool("enablewindowskey")
	d.SeekRel(4)
	d.FieldBool("remoteconsoleaudio")
	d.FieldBool("force_encrypted_cs_pdu")
	d.FieldBool("rail")
	d.FieldBool("logonerrors")
	d.FieldBool("mouse_has_wheel")
	d.FieldBool("password_is_sc_pin")
	d.FieldBool("noaudioplayback")
	d.FieldBool("using_saved_creds")
	d.FieldBool("audiocapture")
	d.FieldBool("video_disable")
	d.FieldBool("reserved1")
	d.FieldBool("reserved2")
	d.FieldBool("hidef_rail_supported")

	d.SeekRel(int64(d.Pos()) % 31)
}
