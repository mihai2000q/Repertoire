import Playlist from 'src/types/models/Playlist.ts'
import { emptySong, reduxRouterRender, withToastify } from '../../test-utils.tsx'
import PlaylistHeaderCard from './PlaylistHeaderCard.tsx'
import userEvent from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'

describe('Playlist Header Card', () => {
  const playlist: Playlist = {
    id: '1',
    title: 'Playlist 1',
    description: "This is the playlist's description",
    createdAt: '',
    updatedAt: '',
    songs: [
      {
        ...emptySong,
        id: '1',
        title: 'Song 1'
      },
      {
        ...emptySong,
        id: '2',
        title: 'Song 2'
      }
    ]
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display minimal info', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<PlaylistHeaderCard playlist={playlist} />)

    expect(screen.getByRole('img', { name: playlist.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: playlist.title })).toHaveAttribute(
      'src',
      playlist.imageUrl
    )
    expect(screen.getByRole('heading', { name: playlist.title })).toBeInTheDocument()
    expect(screen.getByText(playlist.description)).toBeInTheDocument()
    expect(screen.getByText(`${playlist.songs.length} songs`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-header' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /info/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should display image modal, when clicking the image', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<PlaylistHeaderCard playlist={{ ...playlist, imageUrl: 'something.png' }} />)

    await user.click(screen.getByRole('img', { name: playlist.title }))
    expect(
      await screen.findByRole('dialog', { name: playlist.title + '-image' })
    ).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display info modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<PlaylistHeaderCard playlist={playlist} />)

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /info/i }))

      expect(await screen.findByRole('dialog', { name: /playlist info/i })).toBeInTheDocument()
    })

    it('should display edit header modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<PlaylistHeaderCard playlist={playlist} />)

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /edit/i }))

      expect(
        await screen.findByRole('dialog', { name: /edit playlist header/i })
      ).toBeInTheDocument()
    })

    it('should display warning modal and delete playlist', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/playlists/${playlist.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRouterRender(withToastify(<PlaylistHeaderCard playlist={playlist} />))

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete playlist/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete playlist/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(window.location.pathname).toBe('/playlists')
      expect(screen.getByText(`${playlist.title} deleted!`)).toBeInTheDocument()
    })
  })

  it('should display edit header modal from edit button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<PlaylistHeaderCard playlist={playlist} />)

    await user.click(screen.getByRole('button', { name: 'edit-header' }))

    expect(await screen.findByRole('dialog', { name: /edit playlist header/i })).toBeInTheDocument()
  })
})
