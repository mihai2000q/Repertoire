import Playlist from 'src/types/models/Playlist.ts'
import { reduxRouterRender, withToastify } from '../../test-utils.tsx'
import PlaylistCard from './PlaylistCard.tsx'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'

describe('Playlist Card', () => {
  const playlist: Playlist = {
    id: '1',
    title: 'Playlist 1',
    description: '',
    createdAt: '',
    updatedAt: '',
    songs: []
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render minimal info', () => {
    reduxRouterRender(<PlaylistCard playlist={playlist} />)

    expect(screen.getByRole('img', { name: playlist.title })).toBeInTheDocument()
    expect(screen.getByText(playlist.title)).toBeInTheDocument()
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<PlaylistCard playlist={playlist} />)

    await user.pointer({ keys: '[MouseRight>]', target: screen.getByRole('img', { name: playlist.title }) })
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display warning modal when clicking on delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/playlists/${playlist.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRouterRender(withToastify(<PlaylistCard playlist={playlist} />))

      await user.pointer({ keys: '[MouseRight>]', target: screen.getByRole('img', { name: playlist.title }) })
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(screen.getByRole('dialog', { name: /delete/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(screen.getByText(`${playlist.title} deleted!`)).toBeInTheDocument()
    })
  })

  it('should navigate on click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<PlaylistCard playlist={playlist} />)

    await user.click(screen.getByRole('img', { name: playlist.title }))
    expect(window.location.pathname).toBe(`/playlist/${playlist.id}`)
  })
})
