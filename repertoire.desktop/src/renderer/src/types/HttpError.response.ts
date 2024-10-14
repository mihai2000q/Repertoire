export default interface HttpErrorResponse {
  data: Error
  status: number
}

type Error = {
  title: string
  status: number
  errors: object
}
