import { emptyArtist, reduxRouterRender } from '../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { afterEach, expect } from 'vitest'
import { RootState } from '../../../state/store.ts'
import Artist from '../../../types/models/Artist.ts'
import HomeRecentlyPlayedArtistCard from './HomeRecentlyPlayedArtistCard.tsx'
import dayjs from 'dayjs'
import { userEvent } from '@testing-library/user-event'

describe('Home Recently Played Artist Card', () => {
  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    imageUrl: 'something.png',
    name: 'Artist 1',
    lastTimePlayed: '2024-12-14T10:30',
    progress: 500
  }

  afterEach(() => (window.location.pathname = '/'))

  it('should render', () => {
    const localArtist = {
      ...artist,
      imageUrl: null
    }

    reduxRouterRender(<HomeRecentlyPlayedArtistCard artist={localArtist} />)

    expect(screen.getByText(localArtist.name)).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    expect(screen.getByLabelText(`default-icon-${localArtist.name}`)).toBeInTheDocument()
    expect(screen.getByText(dayjs(localArtist.lastTimePlayed).format('DD MMM'))).toBeInTheDocument()
  })

  it('should render with image', () => {
    reduxRouterRender(<HomeRecentlyPlayedArtistCard artist={artist} />)

    expect(screen.getByRole('img', { name: artist.name })).toHaveAttribute('src', artist.imageUrl)
  })

  it('should open artist drawer on clicking', async () => {
    const user = userEvent.setup()

    const [_, store] = reduxRouterRender(<HomeRecentlyPlayedArtistCard artist={artist} />)

    await user.click(await screen.findByText(artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(artist.id)
  })

  it('should display menu on artist right click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<HomeRecentlyPlayedArtistCard artist={artist} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: await screen.findByRole('img', { name: artist.name })
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to artist when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<HomeRecentlyPlayedArtistCard artist={artist} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: await screen.findByRole('img', { name: artist.name })
      })
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/artist/${artist.id}`)
    })
  })
})
