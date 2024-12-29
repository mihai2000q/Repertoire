import { mantineRender } from '../../../test-utils.tsx'
import RightSideEntityDrawer from './RightSideEntityDrawer.tsx'
import { screen } from '@testing-library/react'

describe('Right Side Entity Drawer', () => {
  it('should display children when opened and not loading', () => {
    // Arrange
    const childrenTestId = 'children-mock'
    const children = <div data-testid={childrenTestId}>Children</div>

    // Act
    mantineRender(
      <RightSideEntityDrawer opened={true} onClose={() => {}} isLoading={false} loader={<></>}>
        {children}
      </RightSideEntityDrawer>
    )

    // Assert
    expect(screen.getByTestId(childrenTestId)).toBeInTheDocument()
  })

  it('should not display children when not opened', () => {
    // Arrange
    const childrenTestId = 'children-mock'
    const children = <div data-testid={childrenTestId}>Children</div>

    // Act
    mantineRender(
      <RightSideEntityDrawer opened={false} onClose={() => {}} isLoading={false} loader={<></>}>
        {children}
      </RightSideEntityDrawer>
    )

    // Assert
    expect(screen.queryByTestId(childrenTestId)).not.toBeInTheDocument()
  })

  it('should not display children when loading and display loader', () => {
    // Arrange
    const childrenTestId = 'children-mock'
    const children = <div data-testid={childrenTestId}>Children</div>

    const loaderTestId = 'loader-mock'
    const loader = <div data-testid={loaderTestId}>Loader</div>

    // Act
    mantineRender(
      <RightSideEntityDrawer opened={true} onClose={() => {}} isLoading={true} loader={loader}>
        {children}
      </RightSideEntityDrawer>
    )

    // Assert
    expect(screen.queryByTestId(childrenTestId)).not.toBeInTheDocument()
    expect(screen.queryByTestId(loaderTestId)).toBeVisible()
  })
})
