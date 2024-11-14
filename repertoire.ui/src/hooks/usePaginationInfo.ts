export default function usePaginationInfo(
  totalCount: number | undefined,
  pageSize: number,
  currentPage: number
): {
  startCount: number,
  endCount: number,
  totalPages: number
} {
  const totalPages = totalCount ? Math.ceil(totalCount / pageSize) : 0

  const startCount = totalCount === 0
    ? 0
    : currentPage === 1
      ? 1
      : pageSize * (currentPage - 1)

  const endCount = totalCount
    ? totalCount === 0
      ? 0
      : currentPage === totalPages
        ? totalCount
        : currentPage * pageSize
    : 0

  return {startCount, endCount, totalPages}
}
