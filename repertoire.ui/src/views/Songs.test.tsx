import { screen } from '@testing-library/react'
import Songs from './Songs.tsx'
import { reduxRender } from '../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import Song from '../types/models/Song.ts'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { userEvent } from '@testing-library/user-event'

describe.skip('Songs', () => {
  const songs: Song[] = [
    {
      id: '1',
      title: 'All for justice',
      description: '',
      isRecorded: false,
      sections: []
    },
    {
      id: '2',
      description: '',
      title: 'Seek and Destroy',
      isRecorded: true,
      sections: []
    }
  ]

  const handlers = [
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = {
        models: songs,
        totalCount: 3
      }
      return HttpResponse.json(response)
    }),
    http.get('/songs/1', async () => {
      return HttpResponse.json(songs[0])
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display relevant info when there are songs', async () => {
    const [{ container }] = reduxRender(<Songs />)

    expect(screen.getByRole('heading', { name: /songs/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /new song/i })).toBeInTheDocument()
    expect(screen.getByTestId('songs-loader')).toBeInTheDocument()
    expect(container.querySelector('.mantine-Loader-root')).toBeInTheDocument()

    expect(await screen.findByTestId('new-song-card')).toBeInTheDocument()
    expect(screen.getAllByTestId(/song-card-/)).toHaveLength(songs.length)
    songs.forEach((song) => expect(screen.getByTestId(`song-card-${song.id}`)).toBeInTheDocument())
    expect(screen.getByTestId('songs-pagination')).toHaveTextContent('1')
  })

  it('should open the add new song modal when clicking the new song button', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(<Songs />)

    const newSongButton = screen.getByRole('button', { name: /new song/i })
    await user.click(newSongButton)

    // Assert
    expect(await screen.findByRole('heading', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should open the add new song modal when clicking the new song card button', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(<Songs />)

    const newSongCardButton = await screen.findByTestId('new-song-card')
    await user.click(newSongCardButton)

    // Assert
    expect(await screen.findByRole('heading', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should open song drawer when clicking on song', async () => {
    // Arrange
    const song = songs[0]
    const user = userEvent.setup()

    // Act
    reduxRender(<Songs />)

    await user.click(await screen.findByTestId(`song-card-${song.id}`))

    // Assert
    expect(await screen.findByRole('heading', { name: song.title })).toBeInTheDocument()
  })

  it('should display not display some info when there are no songs', async () => {
    reduxRender(<Songs />)

    server.use(
      http.get('/songs', async () => {
        const response: WithTotalCountResponse<Song> = { models: [], totalCount: 0 }
        return HttpResponse.json(response)
      })
    )

    expect(await screen.findByText(/no songs/)).toBeInTheDocument()
    expect(screen.queryByTestId('new-song-card')).not.toBeInTheDocument()
    expect(screen.queryAllByTestId(/song-card-/)).toHaveLength(0)
    expect(screen.queryByTestId('songs-pagination')).not.toBeInTheDocument()
  })
})
