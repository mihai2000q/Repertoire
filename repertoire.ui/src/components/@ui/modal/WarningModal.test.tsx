import { mantineRender } from '../../../test-utils.tsx'
import WarningModal from './WarningModal.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

describe('Warning Modal', () => {
  it('should render and display content - when description is string', () => {
    const title = 'Warning Modal'
    const description = 'Warning Modal description'

    mantineRender(
      <WarningModal
        opened={true}
        onClose={() => {}}
        title={title}
        description={description}
        onYes={() => {}}
      />
    )

    expect(screen.getByRole('dialog', { name: title })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: title })).toBeInTheDocument()
    expect(screen.getByText(description)).toBeInTheDocument()
  })

  it('should render and display content - when description is React Node', () => {
    const title = 'Warning Modal'

    const descriptionTestId = 'description-mock'
    const description = <div data-testid={descriptionTestId}>Warning Modal Description</div>

    mantineRender(
      <WarningModal
        opened={true}
        onClose={() => {}}
        title={title}
        description={description}
        onYes={() => {}}
      />
    )

    expect(screen.getByRole('dialog', { name: title })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: title })).toBeInTheDocument()
    expect(screen.getByTestId(descriptionTestId)).toBeInTheDocument()
  })

  it('should close modal when clicking Cancel button', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()

    mantineRender(
      <WarningModal
        opened={true}
        onClose={onClose}
        title={'Title'}
        description={'Description'}
        onYes={() => {}}
      />
    )

    const cancelButton = screen.getByRole('button', { name: /cancel/i })
    await user.click(cancelButton)

    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should close modal and invoke onYes when clicking Yes button', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()
    const onYes = vitest.fn()

    mantineRender(
      <WarningModal
        opened={true}
        onClose={onClose}
        title={'Title'}
        description={'Description'}
        onYes={onYes}
      />
    )

    const yesButton = screen.getByRole('button', { name: /yes/i })
    await user.click(yesButton)

    expect(onClose).toHaveBeenCalledOnce()
    expect(onYes).toHaveBeenCalledOnce()
  })
})
