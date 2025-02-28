import SearchType from "../../utils/enums/SearchType.ts";

export interface GetSearchRequest {
  query: string
  currentPage?: number
  pageSize?: number
  type?: SearchType
}
