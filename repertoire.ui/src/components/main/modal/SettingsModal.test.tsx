import { emptyUser, reduxRouterRender } from '../../../test-utils.tsx'
import SettingsModal from './SettingsModal.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

describe('Settings Modal', () => {
  it('should render', () => {
    reduxRouterRender(<SettingsModal opened={true} onClose={() => {}} user={emptyUser} />)

    expect(screen.getByRole('dialog', { name: /settings/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /settings/i })).toBeInTheDocument()
    expect(screen.getByRole('tablist')).toBeInTheDocument()

    expect(screen.getByRole('tab', { name: /account/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /customization/i })).toBeInTheDocument()

    expect(screen.getByRole('tabpanel', { name: /account/i })).toBeInTheDocument()
  })

  it('should be bale to change tabs', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<SettingsModal opened={true} onClose={() => {}} user={emptyUser} />)

    await user.click(screen.getByRole('tab', { name: /customization/i }))
    expect(screen.getByRole('tabpanel', { name: /customization/i })).toBeInTheDocument()

    await user.click(screen.getByRole('tab', { name: /account/i }))
    expect(screen.getByRole('tabpanel', { name: /account/i })).toBeInTheDocument()
  })
})
