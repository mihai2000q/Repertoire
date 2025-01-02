import { renderHook } from '@testing-library/react'
import usePaginationInfo from './usePaginationInfo.ts'

describe('use Pagination Info', () => {
  it.each([
    [
      { totalCount: 21, pageSize: 10, currentPage: 1 },
      { startCount: 1, endCount: 10, totalPages: 3 }
    ],
    [
      { totalCount: 21, pageSize: 10, currentPage: 2 },
      { startCount: 11, endCount: 20, totalPages: 3 }
    ],
    [
      { totalCount: 21, pageSize: 10, currentPage: 3 },
      { startCount: 21, endCount: 21, totalPages: 3 }
    ],
    [
      { totalCount: undefined, pageSize: 10, currentPage: 3 },
      { startCount: 0, endCount: 0, totalPages: 0 }
    ],
    [
      { totalCount: 0, pageSize: 10, currentPage: 3 },
      { startCount: 0, endCount: 0, totalPages: 0 }
    ]
  ])(
    'should return pagination info based on total count, page size and current page',
    (input, output) => {
      // Arrange - parameterized
      // Act
      const { result } = renderHook(() =>
        usePaginationInfo(input.totalCount, input.pageSize, input.currentPage)
      )

      // Assert
      expect(result.current).toStrictEqual(output)
    }
  )
})
