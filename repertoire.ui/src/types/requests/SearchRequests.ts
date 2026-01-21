import SearchType from '../enums/SearchType.ts'

export interface GetSearchRequest {
  query: string
  currentPage?: number
  pageSize?: number
  type?: SearchType
  filter?: string[]
  order?: string[]
  ids?: string[]
}
