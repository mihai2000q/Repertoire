import { emptyUser, reduxRouterRender } from '../../../test-utils.tsx'
import SettingsModalAccountTab from './SettingsModalAccountTab.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

describe('Settings Modal Account Tab', () => {
  it('should render', () => {
    reduxRouterRender(<SettingsModalAccountTab onCloseSettingsModal={() => {}} user={emptyUser} />)

    expect(screen.getByText('Email')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /change email/i })).toBeInTheDocument()
    expect(screen.getByText('Password')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /change password/i })).toBeInTheDocument()
    expect(screen.getByText('Deletion')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /delete account/i })).toBeInTheDocument()
  })

  it('should open delete account modal when clicking on delete account', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<SettingsModalAccountTab onCloseSettingsModal={() => {}} user={emptyUser} />)

    await user.click(screen.getByRole('button', { name: /delete account/i }))
    expect(screen.getByRole('dialog', { name: /delete account/i })).toBeInTheDocument()
  })
})
