import { screen } from '@testing-library/react'
import Songs from './Songs'
import { reduxRender } from '../test-utils'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import Song from "../types/models/Song.ts";
import WithTotalCountResponse from "../types/responses/WithTotalCountResponse.ts";
import {userEvent} from "@testing-library/user-event";

describe('Songs', () => {
  const songs: Song[] = [
    {
      id: '1',
      title: 'All for justice',
      isRecorded: false
    },
    {
      id: '2',
      title: 'Seek and Destroy',
      isRecorded: true
    }
  ]

  const handlers = [
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = {
        models: songs,
        totalCount: 3
      }
      return HttpResponse.json(response)
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
    expect(container.querySelector('.mantine-Loader-root')).toBeInTheDocument()

    expect(await screen.findByTestId('new-song-card')).toBeInTheDocument()
    expect(screen.getAllByTestId(/song-card-/)).toHaveLength(songs.length)
    songs.forEach(song => expect(screen.getByTestId(`song-card-${song.id}`)).toBeInTheDocument())
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
    expect(screen.getByRole('heading', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should open the add new song modal when clicking the new song card button', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(<Songs />)

    const newSongCardButton = await screen.findByTestId('new-song-card')
    await user.click(newSongCardButton)

    // Assert
    expect(screen.getByRole('heading', { name: /add new song/i })).toBeInTheDocument()
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
