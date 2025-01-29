import { reduxRender, withToastify } from '../../../test-utils.tsx'
import { setupServer } from 'msw/node'
import AddNewAlbumSongModal from './AddNewAlbumSongModal.tsx'
import { act, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { CreateSongRequest } from '../../../types/requests/SongRequests.ts'
import Album from '../../../types/models/Album.ts'

describe('Add New Album Song Modal', () => {
  const album: Album = {
    id: '1',
    title: '',
    songs: [],
    createdAt: '',
    updatedAt: ''
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    reduxRender(<AddNewAlbumSongModal opened={true} onClose={() => {}} album={undefined} />)

    expect(screen.getByRole('dialog', { name: /add new song/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /add new song/i })).toBeInTheDocument()
    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /title/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /submit/i })).toBeInTheDocument()
  })

  it('should render and display info when the song will inherit values', () => {
    const localAlbum: Album = {
      ...album,
      releaseDate: '2024-12-10',
      imageUrl: 'something.png',
      artist: {
        id: '',
        name: '',
        albums: [],
        songs: [],
        createdAt: '',
        updatedAt: ''
      }
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

  it('should send only create request when no image is uploaded', async () => {
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
})
