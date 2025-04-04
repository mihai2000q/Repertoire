import { screen } from '@testing-library/react'
import { routerRender } from '../../test-utils.tsx'
import UnknownAlbumCard from './UnknownAlbumCard.tsx'
import userEvent from '@testing-library/user-event'

describe('Unknown Album Card', () => {
  it('should render', () => {
    routerRender(<UnknownAlbumCard />)

    expect(screen.getByLabelText('icon-unknown-album')).toBeInTheDocument()
    expect(screen.getByText(/unknown/i)).toBeInTheDocument()
  })

  it('should navigate on click', async () => {
    const user = userEvent.setup()

    routerRender(<UnknownAlbumCard />)

    await user.click(screen.getByLabelText('icon-unknown-album'))

    expect(window.location.pathname).toBe('/album/unknown')
  })
})
