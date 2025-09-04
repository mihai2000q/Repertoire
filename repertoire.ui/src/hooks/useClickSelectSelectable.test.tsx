import { render, renderHook, screen } from '@testing-library/react'
import { afterEach, beforeAll } from 'vitest'
import useClickSelectSelectable from './useClickSelectSelectable.ts'
import { useClickSelect } from '../context/ClickSelectContext.tsx'

describe('use Click Select Selectable', () => {
  const mockAddSelectables = vi.fn()
  const mockRemoveSelectables = vi.fn()
  const dataTestId = 'data-test-id'
  const id = 'something'

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: mockAddSelectables,
      removeSelectable: mockRemoveSelectables,
      selectedIds: [],
      clearSelection: vi.fn()
    })
  })

  beforeAll(() => {
    // Mock the context
    vi.mock('../context/ClickSelectContext.tsx', () => ({
      useClickSelect: vi.fn()
    }))
  })

  afterEach(() => vi.restoreAllMocks())

  const TestComponent = () => {
    const { ref } = useClickSelectSelectable<HTMLDivElement>(id)
    return <div ref={ref} data-testid={dataTestId} />
  }

  it('should return a ref object, is drag selected and is drag selecting with default values', () => {
    const { result } = renderHook(() => useClickSelectSelectable(''))

    expect(result.current.ref).toBeDefined()
    expect(result.current.ref).toStrictEqual({ current: undefined })
    expect(result.current.isClickSelected).toBeFalsy()
    expect(result.current.isClickSelectionActive).toBeFalsy()
    expect(result.current.isLastInSelection).toBeFalsy()
  })

  it('should call addSelectable when ref is set and add it to selectables', () => {
    render(<TestComponent />)

    expect(mockAddSelectables).toHaveBeenCalledExactlyOnceWith(id, screen.getByTestId(dataTestId))
  })

  it('should call removeSelectable on cleanup', () => {
    const { unmount } = render(<TestComponent />)

    expect(mockAddSelectables).toHaveBeenCalledExactlyOnceWith(id, screen.getByTestId(dataTestId))

    unmount()

    expect(mockRemoveSelectables).toHaveBeenCalledExactlyOnceWith(id)
  })

  it('should return is click selection active true when there are ids', () => {
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: mockAddSelectables,
      removeSelectable: mockRemoveSelectables,
      selectedIds: ['something'],
      clearSelection: vi.fn()
    })

    const { result } = renderHook(() => useClickSelectSelectable('asd'))

    expect(result.current.isClickSelected).toBeFalsy()
    expect(result.current.isClickSelectionActive).toBeTruthy()
  })

  it('should return is click selected true when there are ids and the id is part of those', () => {
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: mockAddSelectables,
      removeSelectable: mockRemoveSelectables,
      selectedIds: ['something', id],
      clearSelection: vi.fn()
    })

    const { result } = renderHook(() => useClickSelectSelectable(id))

    expect(result.current.isClickSelected).toBeTruthy()
    expect(result.current.isClickSelectionActive).toBeTruthy()
  })

  it.each([
    ['id', [{ id: 'id', selected: true }], ['id']],
    [
      'id',
      [
        { id: 'something', selected: true },
        { id: 'something2', selected: false },
        { id: 'id', selected: true }
      ],
      ['something', 'id']
    ],
    [
      'id',
      [
        { id: 'something', selected: true },
        { id: 'id', selected: true },
        { id: 'something2', selected: false }
      ],
      ['something', 'id']
    ],
    [
      'id',
      [
        { id: 'something', selected: true },
        { id: 'id', selected: true },
        { id: 'something2', selected: false },
        { id: 'something3', selected: true }
      ],
      ['something', 'id', 'something3']
    ]
  ])(
    'should return is last in selection true when it is the last selected item or last item in a row (before an unselected one)',
    (id, selectables, selectedIds) => {
      vi.mocked(useClickSelect).mockReturnValue({
        selectables: selectables,
        addSelectable: mockAddSelectables,
        removeSelectable: mockRemoveSelectables,
        selectedIds: selectedIds,
        clearSelection: vi.fn()
      })

      const { result } = renderHook(() => useClickSelectSelectable(id))

      expect(result.current.isLastInSelection).toBeTruthy()
    }
  )

  it.each([
    [
      'id',
      [
        { id: 'something', selected: false },
        { id: 'id', selected: true },
        { id: 'something2', selected: true }
      ],
      ['something', 'id']
    ]
  ])('should return is last in selection false', (id, selectables, selectedIds) => {
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: selectables,
      addSelectable: mockAddSelectables,
      removeSelectable: mockRemoveSelectables,
      selectedIds: selectedIds,
      clearSelection: vi.fn()
    })

    const { result } = renderHook(() => useClickSelectSelectable(id))

    expect(result.current.isLastInSelection).toBeFalsy()
  })
})
