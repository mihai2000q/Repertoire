export default interface User {
  id: string
  name: string
  email: string
  createdAt: string,
  updatedAt: string,
  profilePictureUrl?: string | null,
}
