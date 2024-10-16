export default interface HttpErrorResponse {
  data: Error
  status: number
}

type Error = {
  error: string
}
