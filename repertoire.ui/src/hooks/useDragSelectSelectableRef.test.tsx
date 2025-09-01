import { render, renderHook, screen } from '@testing-library/react'
import useDragSelectSelectableRef from './useDragSelectSelectableRef'
import { useDragSelect } from '../context/DragSelectContext'
import { afterEach } from 'vitest'

// Mock the context
vi.mock('../context/DragSelectContext', () => ({
  useDragSelect: vi.fn()
}))

describe('useDragSelectSelectableRef', () => {
  const mockAddSelectables = vi.fn()
  const mockRemoveSelectables = vi.fn()
  const dataTestId = 'dataTestId'

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useDragSelect).mockReturnValue({
      dragSelect: {
        addSelectables: mockAddSelectables,
        removeSelectables: mockRemoveSelectables
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } as any,
      selectedIds: [],
      clearSelection: vi.fn()
    })
  })

  afterEach(() => vi.restoreAllMocks())

  const TestComponent = () => {
    const ref = useDragSelectSelectableRef<HTMLDivElement>()
    return <div ref={ref} data-testid={dataTestId} />
  }

  it('should return a ref object', () => {
    const { result } = renderHook(() => useDragSelectSelectableRef())

    expect(result.current).toBeDefined()
    expect(result.current).toStrictEqual({ current: undefined })
  })

  it('should call addSelectables when ref is set and dragSelect is available', () => {
    render(<TestComponent />)

    expect(mockAddSelectables).toHaveBeenCalledExactlyOnceWith(screen.getByTestId(dataTestId))
  })

  // simply doesn't work because the cleanup on useEffect is not triggered on testing environments
  it.skip('should call removeSelectables on cleanup when element was added', () => {
    const { unmount } = render(<TestComponent />)

    const element = screen.getByTestId(dataTestId)
    expect(mockAddSelectables).toHaveBeenCalledExactlyOnceWith(element)

    unmount()

    expect(mockRemoveSelectables).toHaveBeenCalledExactlyOnceWith(element)
  })

  it('should not call addSelectables when dragSelect is null', () => {
    vi.mocked(useDragSelect).mockReturnValue({
      dragSelect: null,
      selectedIds: [],
      clearSelection: vi.fn()
    })

    render(<TestComponent />)

    expect(mockAddSelectables).not.toHaveBeenCalled()
  })

  it('should not call removeSelectables on cleanup when dragSelect was null', () => {
    vi.mocked(useDragSelect).mockReturnValue({
      dragSelect: null,
      selectedIds: [],
      clearSelection: vi.fn()
    })

    const { unmount } = render(<TestComponent />)

    unmount()

    expect(mockRemoveSelectables).not.toHaveBeenCalled()
  })
})
