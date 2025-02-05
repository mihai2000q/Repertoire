import { FileWithPath } from '@mantine/dropzone'

export interface UpdateUserRequest {
  name: string
}

export interface SaveProfilePictureRequest {
  profile_pic: FileWithPath
}
