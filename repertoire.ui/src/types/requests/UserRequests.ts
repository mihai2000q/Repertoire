import { FileWithPath } from '@mantine/dropzone'

export interface SignUpRequest {
  name: string
  email: string
  password: string
}

export interface UpdateUserRequest {
  name: string
}

export interface SaveProfilePictureRequest {
  profile_pic: FileWithPath
}
