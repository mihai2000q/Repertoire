import { screen } from '@testing-library/react'
import { routerRender } from '../../test-utils.tsx'
import UnknownArtistCard from './UnknownArtistCard.tsx'
import userEvent from '@testing-library/user-event'

describe('Unknown Artist Card', () => {
  it('should render', () => {
    routerRender(<UnknownArtistCard />)

    expect(screen.getByRole('img', { name: 'unknown-artist' })).toBeInTheDocument()
    expect(screen.getByText(/unknown/i)).toBeInTheDocument()
  })

  it('should navigate on click', async () => {
    const user = userEvent.setup()

    routerRender(<UnknownArtistCard />)

    await user.click(screen.getByRole('img', { name: 'unknown-artist' }))

    expect(window.location.pathname).toBe('/artist/unknown')
  })
})
