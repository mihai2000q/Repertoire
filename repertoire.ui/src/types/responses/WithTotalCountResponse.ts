export default interface WithTotalCountResponse<TModel> {
  models: TModel[]
  totalCount: number
}
