import { FileWithPath } from '@mantine/dropzone'

export interface UpdateUserRequest {
  name: string
}

export interface SaveProfilePictureToUserRequest {
  profile_pic: FileWithPath
}
