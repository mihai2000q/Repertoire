import { reduxRender, withToastify } from '../../../test-utils.tsx'
import { setupServer } from 'msw/node'
import AddNewArtistAlbumModal from './AddNewArtistAlbumModal.tsx'
import { act, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { CreateAlbumRequest } from '../../../types/requests/AlbumRequests.ts'
import dayjs from 'dayjs'

describe('Add New Artist Album Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    reduxRender(<AddNewArtistAlbumModal opened={true} onClose={() => {}} artistId={''} />)

    expect(screen.getByRole('dialog', { name: /add new album/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /add new album/i })).toBeInTheDocument()
    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /title/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /submit/i })).toBeInTheDocument()
  })

  describe('should send only create request when no image is uploaded', () => {
    it('minimal', async () => {
      const user = userEvent.setup()

      const newTitle = 'New Album'

      const onClose = vitest.fn()

      let capturedRequest: CreateAlbumRequest
      server.use(
        http.post('/albums', async (req) => {
          capturedRequest = (await req.request.json()) as CreateAlbumRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(
        withToastify(
          <AddNewArtistAlbumModal opened={true} onClose={onClose} artistId={undefined} />
        )
      )

      await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
      await user.click(screen.getByRole('button', { name: /submit/i }))

      await waitFor(() =>
        expect(capturedRequest).toStrictEqual({
          title: newTitle
        })
      )
      expect(onClose).toHaveBeenCalledOnce()

      expect(screen.getByText(`${newTitle} added!`))
      expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue('')
    })

    it('with release date', async () => {
      const user = userEvent.setup()

      const newTitle = 'New Album'
      const now = new Date()
      const newReleaseDate = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0, 0)

      const onClose = vitest.fn()

      let capturedRequest: CreateAlbumRequest
      server.use(
        http.post('/albums', async (req) => {
          capturedRequest = (await req.request.json()) as CreateAlbumRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(
        withToastify(
          <AddNewArtistAlbumModal opened={true} onClose={onClose} artistId={undefined} />
        )
      )

      await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
      await user.type(screen.getByRole('button', { name: /release-date/i }), newTitle)
      await user.click(
        screen.getByRole('button', { name: dayjs(newReleaseDate).format('D MMMM YYYY') })
      )
      await user.click(screen.getByRole('button', { name: /submit/i }))

      await waitFor(() =>
        expect(capturedRequest).toStrictEqual({
          title: newTitle,
          releaseDate: newReleaseDate.toISOString()
        })
      )
      expect(onClose).toHaveBeenCalledOnce()

      expect(screen.getByText(`${newTitle} added!`))
      expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue('')
    })
  })

  it('should send only create request when no image is uploaded', async () => {
    const user = userEvent.setup()

    const artistId = 'some-artist-id'
    const newTitle = 'New Album'

    const onClose = vitest.fn()

    let capturedRequest: CreateAlbumRequest
    server.use(
      http.post('/albums', async (req) => {
        capturedRequest = (await req.request.json()) as CreateAlbumRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<AddNewArtistAlbumModal opened={true} onClose={onClose} artistId={artistId} />)
    )

    await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
    await user.click(screen.getByRole('button', { name: /submit/i }))

    await waitFor(() =>
      expect(capturedRequest).toStrictEqual({
        title: newTitle,
        artistId: artistId
      })
    )
    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newTitle} added!`))
    expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue('')
  })

  it('should send create request and save image request when the image is uploaded', async () => {
    const user = userEvent.setup()

    const artistId = 'some-artist-id'
    const newTitle = 'New Album'
    const newImage = new File(['something'], 'image.png', { type: 'image/png' })

    const onClose = vitest.fn()

    const returnedId = 'the-album-id'

    let capturedCreateRequest: CreateAlbumRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.post('/albums', async (req) => {
        capturedCreateRequest = (await req.request.json()) as CreateAlbumRequest
        return HttpResponse.json({ id: returnedId })
      }),
      http.put('/albums/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<AddNewArtistAlbumModal opened={true} onClose={onClose} artistId={artistId} />)
    )

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
    await user.click(screen.getByRole('button', { name: /submit/i }))

    await waitFor(() =>
      expect(capturedCreateRequest).toStrictEqual({
        title: newTitle,
        artistId: artistId
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

    reduxRender(<AddNewArtistAlbumModal opened={true} onClose={onClose} artistId={''} />)

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()

    await user.click(screen.getAllByRole('button').find((b) => b.className.includes('CloseButton')))

    expect(await screen.findByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should validate title', async () => {
    const user = userEvent.setup()

    reduxRender(<AddNewArtistAlbumModal opened={true} onClose={() => {}} artistId={''} />)

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
