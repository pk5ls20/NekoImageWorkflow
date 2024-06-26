// Code generated by Fastpb v0.0.2. DO NOT EDIT.

package clientTransform

import (
	fmt "fmt"
	fastpb "github.com/cloudwego/fastpb"
)

var (
	_ = fmt.Errorf
	_ = fastpb.Skip
)

func (x *ClientInfo) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_ClientInfo[number], err)
}

func (x *ClientInfo) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.ClientUUID, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *ClientInfo) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.ClientName, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *PreUploadFileData) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_PreUploadFileData[number], err)
}

func (x *PreUploadFileData) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	var v int32
	v, offset, err = fastpb.ReadInt32(buf, _type)
	if err != nil {
		return offset, err
	}
	x.ScraperType = ScraperType(v)
	return offset, nil
}

func (x *PreUploadFileData) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.ResourceUUID, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *PreUploadFileData) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.ResourceUri, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *UploadFileData) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_UploadFileData[number], err)
}

func (x *UploadFileData) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	var v int32
	v, offset, err = fastpb.ReadInt32(buf, _type)
	if err != nil {
		return offset, err
	}
	x.ScraperType = ScraperType(v)
	return offset, nil
}

func (x *UploadFileData) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.FileUUID, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *UploadFileData) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.FileContent, offset, err = fastpb.ReadBytes(buf, _type)
	return offset, err
}

func (x *FilePreRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FilePreRequest[number], err)
}

func (x *FilePreRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	var v ClientInfo
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.ClientInfo = &v
	return offset, nil
}

func (x *FilePreRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	var v PreUploadFileData
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.Data = append(x.Data, &v)
	return offset, nil
}

func (x *FilePostRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FilePostRequest[number], err)
}

func (x *FilePostRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	var v ClientInfo
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.ClientInfo = &v
	return offset, nil
}

func (x *FilePostRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	var v UploadFileData
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.Data = append(x.Data, &v)
	return offset, nil
}

func (x *FilePreStatusData) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FilePreStatusData[number], err)
}

func (x *FilePreStatusData) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.ResourceUUID, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *FilePreStatusData) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	var v int32
	v, offset, err = fastpb.ReadInt32(buf, _type)
	if err != nil {
		return offset, err
	}
	x.FilePreUploadStatus = FilePreStatusCode(v)
	return offset, nil
}

func (x *FilePostStatusData) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FilePostStatusData[number], err)
}

func (x *FilePostStatusData) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.FileUUID, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *FilePostStatusData) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	var v int32
	v, offset, err = fastpb.ReadInt32(buf, _type)
	if err != nil {
		return offset, err
	}
	x.FilePostUploadStatus = FilePostStatusCode(v)
	return offset, nil
}

func (x *FilePreResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FilePreResponse[number], err)
}

func (x *FilePreResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	var v int32
	v, offset, err = fastpb.ReadInt32(buf, _type)
	if err != nil {
		return offset, err
	}
	x.StatusCode = ResponseStatusCode(v)
	return offset, nil
}

func (x *FilePreResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	var v FilePreStatusData
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.FilePreUploadStatus = append(x.FilePreUploadStatus, &v)
	return offset, nil
}

func (x *FilePreResponse) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.Message, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *FilePostResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FilePostResponse[number], err)
}

func (x *FilePostResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	var v int32
	v, offset, err = fastpb.ReadInt32(buf, _type)
	if err != nil {
		return offset, err
	}
	x.StatusCode = ResponseStatusCode(v)
	return offset, nil
}

func (x *FilePostResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	var v FilePostStatusData
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.FilePostUploadStatus = append(x.FilePostUploadStatus, &v)
	return offset, nil
}

func (x *FilePostResponse) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.Message, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *ClientInfo) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *ClientInfo) fastWriteField1(buf []byte) (offset int) {
	if x.ClientUUID == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 1, x.GetClientUUID())
	return offset
}

func (x *ClientInfo) fastWriteField2(buf []byte) (offset int) {
	if x.ClientName == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetClientName())
	return offset
}

func (x *PreUploadFileData) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *PreUploadFileData) fastWriteField1(buf []byte) (offset int) {
	if x.ScraperType == 0 {
		return offset
	}
	offset += fastpb.WriteInt32(buf[offset:], 1, int32(x.GetScraperType()))
	return offset
}

func (x *PreUploadFileData) fastWriteField2(buf []byte) (offset int) {
	if x.ResourceUUID == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetResourceUUID())
	return offset
}

func (x *PreUploadFileData) fastWriteField3(buf []byte) (offset int) {
	if x.ResourceUri == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 3, x.GetResourceUri())
	return offset
}

func (x *UploadFileData) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *UploadFileData) fastWriteField1(buf []byte) (offset int) {
	if x.ScraperType == 0 {
		return offset
	}
	offset += fastpb.WriteInt32(buf[offset:], 1, int32(x.GetScraperType()))
	return offset
}

func (x *UploadFileData) fastWriteField2(buf []byte) (offset int) {
	if x.FileUUID == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetFileUUID())
	return offset
}

func (x *UploadFileData) fastWriteField3(buf []byte) (offset int) {
	if len(x.FileContent) == 0 {
		return offset
	}
	offset += fastpb.WriteBytes(buf[offset:], 3, x.GetFileContent())
	return offset
}

func (x *FilePreRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *FilePreRequest) fastWriteField1(buf []byte) (offset int) {
	if x.ClientInfo == nil {
		return offset
	}
	offset += fastpb.WriteMessage(buf[offset:], 1, x.GetClientInfo())
	return offset
}

func (x *FilePreRequest) fastWriteField2(buf []byte) (offset int) {
	if x.Data == nil {
		return offset
	}
	for i := range x.GetData() {
		offset += fastpb.WriteMessage(buf[offset:], 2, x.GetData()[i])
	}
	return offset
}

func (x *FilePostRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *FilePostRequest) fastWriteField1(buf []byte) (offset int) {
	if x.ClientInfo == nil {
		return offset
	}
	offset += fastpb.WriteMessage(buf[offset:], 1, x.GetClientInfo())
	return offset
}

func (x *FilePostRequest) fastWriteField2(buf []byte) (offset int) {
	if x.Data == nil {
		return offset
	}
	for i := range x.GetData() {
		offset += fastpb.WriteMessage(buf[offset:], 2, x.GetData()[i])
	}
	return offset
}

func (x *FilePreStatusData) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *FilePreStatusData) fastWriteField1(buf []byte) (offset int) {
	if x.ResourceUUID == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 1, x.GetResourceUUID())
	return offset
}

func (x *FilePreStatusData) fastWriteField2(buf []byte) (offset int) {
	if x.FilePreUploadStatus == 0 {
		return offset
	}
	offset += fastpb.WriteInt32(buf[offset:], 2, int32(x.GetFilePreUploadStatus()))
	return offset
}

func (x *FilePostStatusData) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *FilePostStatusData) fastWriteField1(buf []byte) (offset int) {
	if x.FileUUID == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 1, x.GetFileUUID())
	return offset
}

func (x *FilePostStatusData) fastWriteField2(buf []byte) (offset int) {
	if x.FilePostUploadStatus == 0 {
		return offset
	}
	offset += fastpb.WriteInt32(buf[offset:], 2, int32(x.GetFilePostUploadStatus()))
	return offset
}

func (x *FilePreResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *FilePreResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteInt32(buf[offset:], 1, int32(x.GetStatusCode()))
	return offset
}

func (x *FilePreResponse) fastWriteField2(buf []byte) (offset int) {
	if x.FilePreUploadStatus == nil {
		return offset
	}
	for i := range x.GetFilePreUploadStatus() {
		offset += fastpb.WriteMessage(buf[offset:], 2, x.GetFilePreUploadStatus()[i])
	}
	return offset
}

func (x *FilePreResponse) fastWriteField3(buf []byte) (offset int) {
	if x.Message == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 3, x.GetMessage())
	return offset
}

func (x *FilePostResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *FilePostResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteInt32(buf[offset:], 1, int32(x.GetStatusCode()))
	return offset
}

func (x *FilePostResponse) fastWriteField2(buf []byte) (offset int) {
	if x.FilePostUploadStatus == nil {
		return offset
	}
	for i := range x.GetFilePostUploadStatus() {
		offset += fastpb.WriteMessage(buf[offset:], 2, x.GetFilePostUploadStatus()[i])
	}
	return offset
}

func (x *FilePostResponse) fastWriteField3(buf []byte) (offset int) {
	if x.Message == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 3, x.GetMessage())
	return offset
}

func (x *ClientInfo) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *ClientInfo) sizeField1() (n int) {
	if x.ClientUUID == "" {
		return n
	}
	n += fastpb.SizeString(1, x.GetClientUUID())
	return n
}

func (x *ClientInfo) sizeField2() (n int) {
	if x.ClientName == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetClientName())
	return n
}

func (x *PreUploadFileData) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *PreUploadFileData) sizeField1() (n int) {
	if x.ScraperType == 0 {
		return n
	}
	n += fastpb.SizeInt32(1, int32(x.GetScraperType()))
	return n
}

func (x *PreUploadFileData) sizeField2() (n int) {
	if x.ResourceUUID == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetResourceUUID())
	return n
}

func (x *PreUploadFileData) sizeField3() (n int) {
	if x.ResourceUri == "" {
		return n
	}
	n += fastpb.SizeString(3, x.GetResourceUri())
	return n
}

func (x *UploadFileData) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *UploadFileData) sizeField1() (n int) {
	if x.ScraperType == 0 {
		return n
	}
	n += fastpb.SizeInt32(1, int32(x.GetScraperType()))
	return n
}

func (x *UploadFileData) sizeField2() (n int) {
	if x.FileUUID == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetFileUUID())
	return n
}

func (x *UploadFileData) sizeField3() (n int) {
	if len(x.FileContent) == 0 {
		return n
	}
	n += fastpb.SizeBytes(3, x.GetFileContent())
	return n
}

func (x *FilePreRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *FilePreRequest) sizeField1() (n int) {
	if x.ClientInfo == nil {
		return n
	}
	n += fastpb.SizeMessage(1, x.GetClientInfo())
	return n
}

func (x *FilePreRequest) sizeField2() (n int) {
	if x.Data == nil {
		return n
	}
	for i := range x.GetData() {
		n += fastpb.SizeMessage(2, x.GetData()[i])
	}
	return n
}

func (x *FilePostRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *FilePostRequest) sizeField1() (n int) {
	if x.ClientInfo == nil {
		return n
	}
	n += fastpb.SizeMessage(1, x.GetClientInfo())
	return n
}

func (x *FilePostRequest) sizeField2() (n int) {
	if x.Data == nil {
		return n
	}
	for i := range x.GetData() {
		n += fastpb.SizeMessage(2, x.GetData()[i])
	}
	return n
}

func (x *FilePreStatusData) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *FilePreStatusData) sizeField1() (n int) {
	if x.ResourceUUID == "" {
		return n
	}
	n += fastpb.SizeString(1, x.GetResourceUUID())
	return n
}

func (x *FilePreStatusData) sizeField2() (n int) {
	if x.FilePreUploadStatus == 0 {
		return n
	}
	n += fastpb.SizeInt32(2, int32(x.GetFilePreUploadStatus()))
	return n
}

func (x *FilePostStatusData) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *FilePostStatusData) sizeField1() (n int) {
	if x.FileUUID == "" {
		return n
	}
	n += fastpb.SizeString(1, x.GetFileUUID())
	return n
}

func (x *FilePostStatusData) sizeField2() (n int) {
	if x.FilePostUploadStatus == 0 {
		return n
	}
	n += fastpb.SizeInt32(2, int32(x.GetFilePostUploadStatus()))
	return n
}

func (x *FilePreResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *FilePreResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeInt32(1, int32(x.GetStatusCode()))
	return n
}

func (x *FilePreResponse) sizeField2() (n int) {
	if x.FilePreUploadStatus == nil {
		return n
	}
	for i := range x.GetFilePreUploadStatus() {
		n += fastpb.SizeMessage(2, x.GetFilePreUploadStatus()[i])
	}
	return n
}

func (x *FilePreResponse) sizeField3() (n int) {
	if x.Message == "" {
		return n
	}
	n += fastpb.SizeString(3, x.GetMessage())
	return n
}

func (x *FilePostResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *FilePostResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeInt32(1, int32(x.GetStatusCode()))
	return n
}

func (x *FilePostResponse) sizeField2() (n int) {
	if x.FilePostUploadStatus == nil {
		return n
	}
	for i := range x.GetFilePostUploadStatus() {
		n += fastpb.SizeMessage(2, x.GetFilePostUploadStatus()[i])
	}
	return n
}

func (x *FilePostResponse) sizeField3() (n int) {
	if x.Message == "" {
		return n
	}
	n += fastpb.SizeString(3, x.GetMessage())
	return n
}

var fieldIDToName_ClientInfo = map[int32]string{
	1: "ClientUUID",
	2: "ClientName",
}

var fieldIDToName_PreUploadFileData = map[int32]string{
	1: "ScraperType",
	2: "ResourceUUID",
	3: "ResourceUri",
}

var fieldIDToName_UploadFileData = map[int32]string{
	1: "ScraperType",
	2: "FileUUID",
	3: "FileContent",
}

var fieldIDToName_FilePreRequest = map[int32]string{
	1: "ClientInfo",
	2: "Data",
}

var fieldIDToName_FilePostRequest = map[int32]string{
	1: "ClientInfo",
	2: "Data",
}

var fieldIDToName_FilePreStatusData = map[int32]string{
	1: "ResourceUUID",
	2: "FilePreUploadStatus",
}

var fieldIDToName_FilePostStatusData = map[int32]string{
	1: "FileUUID",
	2: "FilePostUploadStatus",
}

var fieldIDToName_FilePreResponse = map[int32]string{
	1: "StatusCode",
	2: "FilePreUploadStatus",
	3: "Message",
}

var fieldIDToName_FilePostResponse = map[int32]string{
	1: "StatusCode",
	2: "FilePostUploadStatus",
	3: "Message",
}
