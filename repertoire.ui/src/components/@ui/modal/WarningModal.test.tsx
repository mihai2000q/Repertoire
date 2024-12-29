import { mantineRender } from '../../../test-utils.tsx'
import WarningModal from './WarningModal.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

describe('Warning Modal', () => {
  it('should render and display content - when description is string', () => {
    // Arrange
    const title = 'Warning Modal'
    const description = 'Warning Modal description'

    // Act
    mantineRender(
      <WarningModal
        opened={true}
        onClose={() => {}}
        title={title}
        description={description}
        onYes={() => {}}
      />
    )

    // Assert
    expect(screen.getByText(title)).toBeInTheDocument()
    expect(screen.getByText(description)).toBeInTheDocument()
  })

  it('should render and display content - when description is React Node', () => {
    // Arrange
    const title = 'Warning Modal'

    const descriptionTestId = 'description-mock'
    const description = <div data-testid={descriptionTestId}>Warning Modal Description</div>

    // Act
    mantineRender(
      <WarningModal
        opened={true}
        onClose={() => {}}
        title={title}
        description={description}
        onYes={() => {}}
      />
    )

    // Assert
    expect(screen.getByText(title)).toBeInTheDocument()
    expect(screen.getByTestId(descriptionTestId)).toBeInTheDocument()
  })

  it('should close modal when clicking Cancel button', async () => {
    // Arrange
    const user = userEvent.setup()

    const onClose = vitest.fn()

    // Act
    mantineRender(
      <WarningModal
        opened={true}
        onClose={onClose}
        title={'Title'}
        description={'Description'}
        onYes={() => {}}
      />
    )

    // Assert
    const cancelButton = screen.getByRole('button', { name: /cancel/i })
    await user.click(cancelButton)

    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should close modal and invoke onYes when clicking Yes button', async () => {
    // Arrange
    const user = userEvent.setup()

    const onClose = vitest.fn()
    const onYes = vitest.fn()

    // Act
    mantineRender(
      <WarningModal
        opened={true}
        onClose={onClose}
        title={'Title'}
        description={'Description'}
        onYes={onYes}
      />
    )

    // Assert
    const yesButton = screen.getByRole('button', { name: /yes/i })
    await user.click(yesButton)

    expect(onClose).toHaveBeenCalledOnce()
    expect(onYes).toHaveBeenCalledOnce()
  })
})
