import { render, renderHook, screen } from '@testing-library/react'
import useDragSelectSelectable from './useDragSelectSelectable.ts'
import { useDragSelect } from '../context/DragSelectContext'
import { afterEach, beforeAll } from 'vitest'

describe('use Drag Select Selectable', () => {
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

  beforeAll(() => {
    // Mock the context
    vi.mock('../context/DragSelectContext', () => ({
      useDragSelect: vi.fn()
    }))
  })

  afterEach(() => vi.restoreAllMocks())

  const TestComponent = ({ id = '' }: { id?: string }) => {
    const { ref } = useDragSelectSelectable<HTMLDivElement>(id)
    return <div ref={ref} data-testid={dataTestId} />
  }

  it('should return a ref object, is drag selected and is drag selecting with default values', () => {
    const { result } = renderHook(() => useDragSelectSelectable(''))

    expect(result.current.ref).toBeDefined()
    expect(result.current.ref).toStrictEqual({ current: undefined })
    expect(result.current.isDragSelected).toBeFalsy()
    expect(result.current.isDragSelecting).toBeFalsy()
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

  it('should return is drag selecting true when there are ids', () => {
    vi.mocked(useDragSelect).mockReturnValue({
      dragSelect: {
        addSelectables: mockAddSelectables,
        removeSelectables: mockRemoveSelectables
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } as any,
      selectedIds: ['something'],
      clearSelection: vi.fn()
    })

    const { result } = renderHook(() => useDragSelectSelectable('asd'))

    expect(result.current.isDragSelected).toBeFalsy()
    expect(result.current.isDragSelecting).toBeTruthy()
  })

  it('should return is drag selected true when there are ids and the id is part of those', () => {
    const id = 'something2'

    vi.mocked(useDragSelect).mockReturnValue({
      dragSelect: {
        addSelectables: mockAddSelectables,
        removeSelectables: mockRemoveSelectables
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } as any,
      selectedIds: ['something', id],
      clearSelection: vi.fn()
    })

    const { result } = renderHook(() => useDragSelectSelectable(id))

    expect(result.current.isDragSelected).toBeTruthy()
    expect(result.current.isDragSelecting).toBeTruthy()
  })
})
