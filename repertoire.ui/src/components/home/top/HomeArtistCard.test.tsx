import { emptyArtist, reduxRouterRender } from '../../../test-utils.tsx'
import HomeArtistCard from './HomeArtistCard.tsx'
import Artist from '../../../types/models/Artist.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../../state/store.ts'
import { afterEach } from 'vitest'

describe('Home Artist Card', () => {
  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1'
  }

  afterEach(() => (window.location.pathname = '/'))

  it('should render with minimal info', () => {
    reduxRouterRender(<HomeArtistCard artist={artist} />)

    expect(screen.getByLabelText(`default-icon-${artist.name}`)).toBeInTheDocument()
    expect(screen.getByText(artist.name)).toBeInTheDocument()
  })

  it('should render with maximal info', () => {
    const localArtist = { ...artist, imageUrl: 'something.png' }

    reduxRouterRender(<HomeArtistCard artist={localArtist} />)

    expect(screen.getByRole('img', { name: localArtist.name })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localArtist.name })).toHaveAttribute(
      'src',
      localArtist.imageUrl
    )
    expect(screen.getByText(localArtist.name)).toBeInTheDocument()
  })

  it('should open artist drawer by clicking on image', async () => {
    const user = userEvent.setup()

    const localArtist = { ...artist, imageUrl: 'something.png' }

    const [_, store] = reduxRouterRender(<HomeArtistCard artist={localArtist} />)

    await user.click(screen.getByRole('img', { name: localArtist.name }))
    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localArtist.id)
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<HomeArtistCard artist={artist} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`default-icon-${artist.name}`)
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to artist when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<HomeArtistCard artist={artist} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${artist.name}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/artist/${artist.id}`)
    })
  })
})
