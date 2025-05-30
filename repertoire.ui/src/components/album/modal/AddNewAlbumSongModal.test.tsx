import { emptyAlbum, emptyArtist, reduxRender, withToastify } from '../../../test-utils.tsx'
import { setupServer } from 'msw/node'
import AddNewAlbumSongModal from './AddNewAlbumSongModal.tsx'
import { act, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { CreateSongRequest } from '../../../types/requests/SongRequests.ts'
import Album from '../../../types/models/Album.ts'
import { GuitarTuning } from '../../../types/models/Song.ts'
import Difficulty from '../../../types/enums/Difficulty.ts'

describe('Add New Album Song Modal', () => {
  const album: Album = {
    ...emptyAlbum,
    id: '1'
  }

  const guitarTunings: GuitarTuning[] = [
    {
      id: '1',
      name: 'E Standard'
    },
    {
      id: '2',
      name: 'Drop D'
    },
    {
      id: '3',
      name: 'Drop A'
    }
  ]

  const handlers = [
    http.get('/songs/guitar-tunings', async () => {
      return HttpResponse.json(guitarTunings)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const validSongsterrLink =
    'https://www.songsterr.com/a/wsa/metallica-master-of-puppets-tab-s455118'
  const validYoutubeLink = 'https://www.youtube.com/watch?v=E0ozmU9cJDg'

  it('should render', () => {
    reduxRender(<AddNewAlbumSongModal opened={true} onClose={() => {}} album={undefined} />)

    expect(screen.getByRole('dialog', { name: /add new song/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /add new song/i })).toBeInTheDocument()
    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /title/i })).toBeInTheDocument()

    expect(screen.getByRole('button', { name: /guitar-tuning/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /difficulty/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /bpm/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /songsterr/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /youtube/i })).toBeInTheDocument()

    expect(screen.getByRole('button', { name: /submit/i })).toBeInTheDocument()
  })

  it('should render and display info when the song will inherit values', () => {
    const localAlbum: Album = {
      ...album,
      releaseDate: '2024-12-10',
      imageUrl: 'something.png',
      artist: emptyArtist
    }

    reduxRender(<AddNewAlbumSongModal opened={true} onClose={() => {}} album={localAlbum} />)

    expect(screen.getByText(/if no image is uploaded, it will be inherited/i)).toBeInTheDocument()
    expect(screen.getByText(/new song will inherit/i)).toBeInTheDocument()
    expect(screen.getByText(/release date/i)).toBeInTheDocument()
    expect(screen.getByText(/artist/i)).toBeInTheDocument()
  })

  it('should send create request even when no album is specified (the album is unknown)', async () => {
    const user = userEvent.setup()

    const newTitle = 'New Song'

    const onClose = vitest.fn()

    let capturedRequest: CreateSongRequest
    server.use(
      http.post('/songs', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<AddNewAlbumSongModal opened={true} onClose={onClose} album={undefined} />)
    )

    await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
    await user.click(screen.getByRole('button', { name: /submit/i }))

    await waitFor(() =>
      expect(capturedRequest).toStrictEqual({
        title: newTitle,
        description: ''
      })
    )
    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newTitle} added!`))
    expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue('')
  })

  describe('should send only create request when no image is uploaded', () => {
    it('minimal', async () => {
      const user = userEvent.setup()

      const newTitle = 'New Song'

      const onClose = vitest.fn()

      let capturedRequest: CreateSongRequest
      server.use(
        http.post('/songs', async (req) => {
          capturedRequest = (await req.request.json()) as CreateSongRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(
        withToastify(<AddNewAlbumSongModal opened={true} onClose={onClose} album={album} />)
      )

      await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
      await user.click(screen.getByRole('button', { name: /submit/i }))

      await waitFor(() =>
        expect(capturedRequest).toStrictEqual({
          title: newTitle,
          description: '',
          albumId: album.id
        })
      )
      expect(onClose).toHaveBeenCalledOnce()

      expect(screen.getByText(`${newTitle} added!`))
      expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue('')
    })

    it('with optional fields', async () => {
      const user = userEvent.setup()

      const newTitle = 'New Song'
      const newGuitarTuning = guitarTunings[1]
      const newDifficulty = Difficulty.Medium
      const newBpm = 123

      const onClose = vitest.fn()

      let capturedRequest: CreateSongRequest
      server.use(
        http.post('/songs', async (req) => {
          capturedRequest = (await req.request.json()) as CreateSongRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(
        withToastify(<AddNewAlbumSongModal opened={true} onClose={onClose} album={album} />)
      )

      // fill form
      await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)

      await user.click(screen.getByRole('button', { name: 'guitar-tuning' }))
      await user.click(screen.getByRole('option', { name: newGuitarTuning.name }))
      expect(screen.getByRole('button', { name: 'guitar-tuning' })).toHaveAttribute(
        'aria-selected',
        'true'
      )

      await user.click(screen.getByRole('button', { name: 'difficulty' }))
      await user.click(screen.getByRole('option', { name: newDifficulty }))
      expect(screen.getByRole('button', { name: 'difficulty' })).toHaveAttribute(
        'aria-selected',
        'true'
      )

      await user.click(screen.getByRole('button', { name: 'bpm' }))
      await user.type(screen.getByRole('textbox', { name: 'bpm' }), newBpm.toString())
      expect(screen.getByRole('button', { name: 'bpm' })).toHaveAttribute('aria-selected', 'true')

      await user.click(screen.getByRole('button', { name: 'songsterr' }))
      await user.click(screen.getByRole('textbox', { name: 'songsterr' }))
      await user.paste(validSongsterrLink)
      expect(screen.getByRole('button', { name: 'songsterr' })).toHaveAttribute(
        'aria-selected',
        'true'
      )

      await user.click(screen.getByRole('button', { name: 'youtube' }))
      await user.click(screen.getByRole('textbox', { name: 'youtube' }))
      await user.paste(validYoutubeLink)
      expect(screen.getByRole('button', { name: 'youtube' })).toHaveAttribute(
        'aria-selected',
        'true'
      )

      await user.click(screen.getByRole('button', { name: /submit/i }))

      await waitFor(() =>
        expect(capturedRequest).toStrictEqual({
          title: newTitle,
          description: '',
          albumId: album.id,
          guitarTuningId: newGuitarTuning.id,
          difficulty: newDifficulty,
          bpm: newBpm,
          songsterrLink: validSongsterrLink,
          youtubeLink: validYoutubeLink
        })
      )
      expect(onClose).toHaveBeenCalledOnce()

      expect(screen.getByText(`${newTitle} added!`))

      // reset form
      expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue('')

      await user.click(screen.getByRole('button', { name: 'guitar-tuning' }))
      expect(screen.getByRole('textbox', { name: 'search' })).toHaveValue('')

      await user.click(screen.getByRole('button', { name: 'difficulty' }))
      expect(screen.getByRole('textbox', { name: 'search' })).toHaveValue('')

      await user.click(screen.getByRole('button', { name: 'bpm' }))
      expect(screen.getByRole('textbox', { name: 'bpm' })).toHaveValue('')

      await user.click(screen.getByRole('button', { name: 'songsterr' }))
      expect(screen.getByRole('textbox', { name: 'songsterr' })).toHaveValue('')

      await user.click(screen.getByRole('button', { name: 'youtube' }))
      expect(screen.getByRole('textbox', { name: 'youtube' })).toHaveValue('')
    })
  })

  it('should send create request and save image request when the image is uploaded', async () => {
    const user = userEvent.setup()

    const newTitle = 'New Song'
    const newImage = new File(['something'], 'image.png', { type: 'image/png' })

    const onClose = vitest.fn()

    const returnedId = 'the-song-id'

    let capturedCreateRequest: CreateSongRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.post('/songs', async (req) => {
        capturedCreateRequest = (await req.request.json()) as CreateSongRequest
        return HttpResponse.json({ id: returnedId })
      }),
      http.put('/songs/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<AddNewAlbumSongModal opened={true} onClose={onClose} album={album} />)
    )

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
    await user.click(screen.getByRole('button', { name: /submit/i }))

    await waitFor(() =>
      expect(capturedCreateRequest).toStrictEqual({
        title: newTitle,
        description: '',
        albumId: album.id
      })
    )
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
    expect(capturedSaveImageFormData.get('id')).toBe(returnedId)

    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newTitle} added!`))
    expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue('')
    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
  })

  it('should remove the image on close', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })

    const onClose = vitest.fn()

    reduxRender(<AddNewAlbumSongModal opened={true} onClose={onClose} album={undefined} />)

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()

    await user.click(screen.getAllByRole('button').find((b) => b.className.includes('CloseButton')))

    expect(await screen.findByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should validate title', async () => {
    const user = userEvent.setup()

    reduxRender(<AddNewAlbumSongModal opened={true} onClose={() => {}} album={undefined} />)

    const title = screen.getByRole('textbox', { name: /title/i })

    expect(title).not.toBeInvalid()
    await user.click(screen.getByRole('button', { name: /submit/i }))
    expect(title).toBeInvalid()

    await user.type(title, 'something')
    expect(title).not.toBeInvalid()

    await user.clear(title)
    act(() => title.blur())
    expect(title).toBeInvalid()
  })

  it("should validate songsterr and youtube links' fields", async () => {
    const user = userEvent.setup()

    reduxRender(<AddNewAlbumSongModal opened={true} onClose={() => {}} album={undefined} />)

    // songsterr
    const songsterrButton = screen.getByRole('button', { name: /songsterr/i })
    await user.click(songsterrButton)
    const songsterr = screen.getByRole('textbox', { name: /songsterr/i })

    expect(songsterrButton).not.toBeInvalid()
    expect(songsterr).not.toBeInvalid()
    await user.type(songsterr, 'something')
    act(() => songsterr.blur())
    expect(songsterr).toBeInvalid()
    expect(songsterrButton).toBeInvalid()

    await user.clear(songsterr)
    act(() => songsterr.blur())
    expect(songsterr).not.toBeInvalid()
    expect(songsterrButton).not.toBeInvalid()

    await user.click(songsterr)
    await user.paste(validSongsterrLink)
    expect(songsterr).not.toBeInvalid()
    expect(songsterrButton).not.toBeInvalid()

    // youtube
    const youtubeButton = screen.getByRole('button', { name: /youtube/i })
    await user.click(youtubeButton)
    const youtube = screen.getByRole('textbox', { name: /youtube/i })

    expect(youtube).not.toBeInvalid()
    expect(youtubeButton).not.toBeInvalid()
    await user.type(youtube, 'something')
    act(() => youtube.blur())
    expect(youtube).toBeInvalid()
    expect(youtubeButton).toBeInvalid()

    await user.clear(youtube)
    act(() => youtube.blur())
    expect(youtube).not.toBeInvalid()
    expect(youtubeButton).not.toBeInvalid()

    await user.click(youtube)
    await user.paste(validYoutubeLink)
    expect(youtube).not.toBeInvalid()
    expect(youtubeButton).not.toBeInvalid()
  })
})
