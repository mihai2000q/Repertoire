import { screen } from '@testing-library/react'
import FiltersDrawer from './FiltersDrawer.tsx'
import Filter from '../../../types/Filter.ts'
import FilterOperator from '../../../types/enums/FilterOperator.ts'
import { mantineRender } from '../../../test-utils.tsx'
import { userEvent } from '@testing-library/user-event'

describe('FiltersDrawer', () => {
  const initialFilters = new Map<string, Filter>([
    ['name=', { property: 'name', operator: FilterOperator.Equal, value: '', isSet: false }]
  ])

  const currentFilters = new Map<string, Filter>([
    ['name=', { property: 'name', operator: FilterOperator.Equal, value: 'John', isSet: true }]
  ])

  const defaultProps = {
    opened: true,
    onClose: vi.fn(),
    initialFilters: initialFilters,
    filters: currentFilters,
    internalFilters: currentFilters,
    setFilters: vi.fn(),
    setInternalFilters: vi.fn()
  }

  it('should render', () => {
    const testId = 'test-content-id'
    const testContent = <div data-testid={testId}>Test Content</div>
    const setInternalFilters = vi.fn()

    mantineRender(
      <FiltersDrawer {...defaultProps} setInternalFilters={setInternalFilters}>
        {testContent}
      </FiltersDrawer>
    )

    expect(screen.getByText('Filters')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /reset/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /apply-filters/i })).toBeInTheDocument()
    expect(screen.getByTestId(testId)).toBeInTheDocument()
    expect(setInternalFilters).toHaveBeenCalledWith(initialFilters)
  })

  describe('Apply Filters button', () => {
    it('should be disabled when internalFilters match filters', () => {
      mantineRender(<FiltersDrawer {...defaultProps} />)
      const applyButton = screen.getByRole('button', { name: /apply-filters/i })
      expect(applyButton).toBeDisabled()
    })

    it('should be enabled when internalFilters differ from filters', () => {
      const differentFilters = new Map<string, Filter>([
        ['name=', { property: 'name', operator: FilterOperator.Equal, value: 'Alice', isSet: true }]
      ])
      mantineRender(<FiltersDrawer {...defaultProps} internalFilters={differentFilters} />)

      expect(screen.getByRole('button', { name: /apply-filters/i })).not.toBeDisabled()
    })

    it('should call setFilters with internalFilters when clicked', async () => {
      const user = userEvent.setup()

      const differentFilters = new Map<string, Filter>([
        ['name=', { property: 'name', operator: FilterOperator.Equal, value: 'Alice', isSet: true }]
      ])
      const setFilters = vi.fn()

      mantineRender(
        <FiltersDrawer
          {...defaultProps}
          internalFilters={differentFilters}
          setFilters={setFilters}
        />
      )

      await user.click(screen.getByRole('button', { name: /apply-filters/i }))

      expect(setFilters).toHaveBeenCalledWith(differentFilters)
    })
  })

  describe('Reset button', () => {
    it('should be disabled when both internalFilters and filters match initialFilters', () => {
      mantineRender(
        <FiltersDrawer
          {...defaultProps}
          filters={initialFilters}
          internalFilters={initialFilters}
        />
      )

      const resetButton = screen.getByRole('button', { name: /reset/i })
      expect(resetButton).toBeDisabled()
    })

    it('should be enabled when filters differ from initialFilters', () => {
      mantineRender(<FiltersDrawer {...defaultProps} />)
      expect(screen.getByRole('button', { name: /reset/i })).not.toBeDisabled()
    })

    it('should call setFilters and setInternalFilters with initialFilters when clicked', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()
      const setInternalFilters = vi.fn()
      const additionalReset = vi.fn()

      mantineRender(
        <FiltersDrawer
          {...defaultProps}
          setFilters={setFilters}
          setInternalFilters={setInternalFilters}
          additionalReset={additionalReset}
        />
      )

      await user.click(screen.getByRole('button', { name: /reset/i }))

      expect(setFilters).toHaveBeenCalledWith(initialFilters)
      expect(setInternalFilters).toHaveBeenCalledWith(initialFilters)
      expect(additionalReset).toHaveBeenCalled()
    })
  })
})
