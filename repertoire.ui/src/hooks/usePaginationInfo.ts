export default function usePaginationInfo(
  totalCount: number | undefined,
  pageSize: number,
  currentPage: number
): {
  startCount: number
  endCount: number
  totalPages: number
} {
  if (!totalCount || totalCount === 0) {
    return { startCount: 0, endCount: 0, totalPages: 0 }
  }

  const totalPages = Math.ceil(totalCount / pageSize)

  const startCount = currentPage === 1 ? 1 : pageSize * (currentPage - 1) + 1

  const endCount = currentPage === totalPages ? totalCount : currentPage * pageSize

  return { startCount, endCount, totalPages }
}
