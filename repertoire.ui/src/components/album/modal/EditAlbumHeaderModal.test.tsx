import { emptyArtist, reduxRender, withToastify } from '../../../test-utils.tsx'
import EditAlbumHeaderModal from './EditAlbumHeaderModal.tsx'
import Album from '../../../types/models/Album.ts'
import { act, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import dayjs from 'dayjs'
import { http, HttpResponse } from 'msw'
import { UpdateAlbumRequest } from '../../../types/requests/AlbumRequests.ts'
import Artist from '../../../types/models/Artist.ts'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'

describe('Edit Album Header Modal', () => {
  const album: Album = {
    id: '1',
    title: 'Album 1',
    songs: [],
    createdAt: '',
    updatedAt: '',
    releaseDate: '2024-12-12T00:00:00',
    imageUrl: 'some-image.png'
  }

  const artists: Artist[] = [
    {
      ...emptyArtist,
      id: '1',
      name: 'Chester'
    },
    {
      ...emptyArtist,
      id: '2',
      name: 'Michael'
    },
    {
      ...emptyArtist,
      id: '3',
      name: 'Luther'
    }
  ]

  const handlers = [
    http.get('/artists', async () => {
      const response: WithTotalCountResponse<Artist> = {
        models: artists,
        totalCount: artists.length
      }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<EditAlbumHeaderModal opened={true} onClose={() => {}} album={album} />)

    expect(screen.getByRole('dialog', { name: /edit album header/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit album header/i })).toBeInTheDocument()

    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: 'image-preview' })).toHaveAttribute(
      'src',
      album.imageUrl
    )

    expect(screen.getByRole('textbox', { name: /title/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /title/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue(album.title)

    expect(screen.getByRole('textbox', { name: /artist/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /artist/i })).toHaveValue(album.artist?.name ?? '')

    expect(screen.getByRole('button', { name: /release date/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /release date/i })).toHaveTextContent(
      dayjs(album.releaseDate).format('D MMMM YYYY')
    )

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should send only edit request when the image is unchanged', async () => {
    const user = userEvent.setup()

    const newTitle = 'New Album'
    const newReleaseDate = dayjs('2024-12-24')
    const newArtist = artists[1]
    const onClose = vitest.fn()

    let capturedRequest: UpdateAlbumRequest
    server.use(
      http.put('/albums', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateAlbumRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditAlbumHeaderModal opened={true} onClose={onClose} album={album} />)
    )

    const titleField = screen.getByRole('textbox', { name: /title/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(titleField)
    await user.type(titleField, newTitle)

    await user.click(screen.getByRole('button', { name: /release date/i }))
    await user.click(screen.getByText(newReleaseDate.date().toString()))

    await user.click(screen.getByRole('textbox', { name: /artist/i }))
    await user.click(screen.getByRole('option', { name: newArtist.name }))

    await user.click(saveButton)

    expect(await screen.findByText(/album updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: album.id,
      title: newTitle,
      releaseDate: newReleaseDate.toISOString(),
      artistId: newArtist.id
    })
  })

  it('should send edit request and save image request when the image is replaced', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedRequest: UpdateAlbumRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/albums', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateAlbumRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.put('/albums/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditAlbumHeaderModal opened={true} onClose={onClose} album={album} />)
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(/album updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: album.id,
      title: album.title,
      releaseDate: dayjs(album.releaseDate).toISOString()
    })
    expect(capturedSaveImageFormData.get('id')).toBe(album.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should send edit request and save image request when the image is first added', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedRequest: UpdateAlbumRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/albums', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateAlbumRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.put('/albums/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(
        <EditAlbumHeaderModal
          opened={true}
          onClose={onClose}
          album={{ ...album, imageUrl: undefined }}
        />
      )
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(/album updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: album.id,
      title: album.title,
      releaseDate: dayjs(album.releaseDate).toISOString()
    })
    expect(capturedSaveImageFormData.get('id')).toBe(album.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should send edit request and delete image request', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()

    let capturedRequest: UpdateAlbumRequest
    server.use(
      http.put('/albums', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateAlbumRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.delete(`/albums/images/${album.id}`, () => {
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditAlbumHeaderModal opened={true} onClose={onClose} album={album} />)
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.click(screen.getByRole('button', { name: 'remove-image' }))
    await user.click(saveButton)

    expect(await screen.findByText(/album updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: album.id,
      title: album.title,
      releaseDate: dayjs(album.releaseDate).toISOString()
    })
  })

  it('should disable the save button when no changes are made', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const day = new Date(album.releaseDate).getDate()
    const newDay = dayjs(album.releaseDate).add(1, 'day').date()
    const newArtist = artists[0]

    reduxRender(<EditAlbumHeaderModal opened={true} onClose={() => {}} album={album} />)

    const titleField = screen.getByRole('textbox', { name: /title/i })
    const releaseDateField = screen.getByRole('button', { name: /release date/i })
    const artistField = screen.getByRole('textbox', { name: /artist/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    // change image
    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset image
    await user.click(screen.getByRole('button', { name: 'reset-image' }))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change title
    await user.type(titleField, '1')
    act(() => titleField.blur())
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset title
    await user.clear(titleField)
    await user.type(titleField, album.title)
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change artist
    await user.click(artistField)
    await user.click(await screen.findByRole('option', { name: newArtist.name }))
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset artist
    await user.click(artistField)
    await user.click(await screen.findByRole('option', { name: newArtist.name }))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change release date
    await user.click(releaseDateField)
    await user.click(screen.getByText(newDay.toString()))
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset release date
    await user.click(releaseDateField)
    await user.click(screen.getByText(day.toString()))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // remove image
    await user.click(screen.getByRole('button', { name: 'remove-image' }))
    expect(saveButton).not.toHaveAttribute('data-disabled')
  })

  it('should display info message when image changes', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })

    reduxRender(<EditAlbumHeaderModal opened={true} onClose={() => {}} album={album} />)

    // change image
    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    expect(await screen.findByText(/update all the associated songs/i)).toBeInTheDocument()

    // reset image
    await user.click(screen.getByRole('button', { name: 'reset-image' }))
    expect(screen.queryByText(/update all the associated songs/i)).not.toBeInTheDocument()
  })

  it('should display info message when artist changes', async () => {
    const user = userEvent.setup()

    const newArtist = artists[0]

    reduxRender(<EditAlbumHeaderModal opened={true} onClose={() => {}} album={album} />)

    // change artist
    await user.click(screen.getByRole('textbox', { name: /artist/i }))
    await user.click(screen.getByRole('option', { name: newArtist.name }))
    expect(await screen.findByText(/update all the associated songs/i)).toBeInTheDocument()

    // reset artist
    await user.click(screen.getByRole('textbox', { name: /artist/i }))
    await user.click(screen.getByRole('option', { name: newArtist.name }))
    await waitFor(() =>
      expect(screen.queryByText(/update all the associated songs/i)).not.toBeInTheDocument()
    )
  })

  it('should validate the title textbox', async () => {
    const user = userEvent.setup()

    reduxRender(<EditAlbumHeaderModal opened={true} onClose={() => {}} album={album} />)

    const titleField = screen.getByRole('textbox', { name: /title/i })
    await user.clear(titleField)
    expect(titleField).toBeInvalid()
  })
})
