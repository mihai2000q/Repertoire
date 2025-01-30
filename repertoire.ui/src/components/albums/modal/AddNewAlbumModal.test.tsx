import {emptyArtist, reduxRender, withToastify} from '../../../test-utils.tsx'
import { setupServer } from 'msw/node'
import AddNewAlbumModal from './AddNewAlbumModal.tsx'
import { act, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { CreateAlbumRequest } from '../../../types/requests/AlbumRequests.ts'
import Artist from '../../../types/models/Artist.ts'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'

describe('Add New Album Modal', () => {
  const artists: Artist[] = [
    {
      ...emptyArtist,
      id: '1',
      name: 'Artist 1',
    },
    {
      ...emptyArtist,
      id: '2',
      name: 'Artist 2',
    },
    {
      ...emptyArtist,
      id: '3',
      name: 'Artist 3',
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

  it('should render', () => {
    reduxRender(<AddNewAlbumModal opened={true} onClose={() => {}} />)

    expect(screen.getByRole('dialog', { name: /add new album/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /add new album/i })).toBeInTheDocument()
    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /title/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /artist/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /release date/i })).toBeInTheDocument()
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

      reduxRender(withToastify(<AddNewAlbumModal opened={true} onClose={onClose} />))

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

    it('with new artist and release date', async () => {
      const user = userEvent.setup()

      const newTitle = 'New Album'
      const newArtistName = 'New Artist Name'
      const now = new Date()
      const newReleaseDate = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0, 0)
      const whichDayIsCurrentDay = newReleaseDate.getDate() > 15 ? 1 : 0

      const onClose = vitest.fn()

      let capturedRequest: CreateAlbumRequest
      server.use(
        http.post('/albums', async (req) => {
          capturedRequest = (await req.request.json()) as CreateAlbumRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(withToastify(<AddNewAlbumModal opened={true} onClose={onClose} />))

      await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
      await user.type(screen.getByRole('textbox', { name: /artist/i }), newArtistName)
      await user.click(screen.getByRole('button', { name: /release date/i }))
      await user.click(screen.getAllByText(newReleaseDate.getDate().toString())[whichDayIsCurrentDay])
      await user.click(screen.getByRole('button', { name: /submit/i }))

      await waitFor(() =>
        expect(capturedRequest).toEqual({
          title: newTitle,
          artistName: newArtistName,
          releaseDate: newReleaseDate.toISOString()
        })
      )
      expect(onClose).toHaveBeenCalledOnce()

      expect(screen.getByText(`${newTitle} added!`))
      expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /artist/i })).toHaveValue('')
      expect(screen.getByRole('button', { name: /release date/i })).toHaveTextContent(/choose/i) // the placeholder
    })

    it('with artist', async () => {
      const user = userEvent.setup()

      const newTitle = 'New Album'
      const artist = artists[1]

      const onClose = vitest.fn()

      let capturedRequest: CreateAlbumRequest
      server.use(
        http.post('/albums', async (req) => {
          capturedRequest = (await req.request.json()) as CreateAlbumRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(withToastify(<AddNewAlbumModal opened={true} onClose={onClose} />))

      await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
      await user.click(screen.getByRole('textbox', { name: /artist/i }))
      await user.click(await screen.findByText(artist.name))
      await user.click(screen.getByRole('button', { name: /submit/i }))

      await waitFor(() =>
        expect(capturedRequest).toEqual({
          title: newTitle,
          artistId: artist.id
        })
      )
      expect(onClose).toHaveBeenCalledOnce()

      expect(screen.getByText(`${newTitle} added!`))
      expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /artist/i })).toHaveValue('')
    })
  })

  it('should send create request and save image request when the image is uploaded', async () => {
    const user = userEvent.setup()

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

    reduxRender(withToastify(<AddNewAlbumModal opened={true} onClose={onClose} />))

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
    await user.click(screen.getByRole('button', { name: /submit/i }))

    await waitFor(() =>
      expect(capturedCreateRequest).toEqual({
        title: newTitle
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

    reduxRender(<AddNewAlbumModal opened={true} onClose={onClose} />)

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()

    await user.click(screen.getAllByRole('button').find((b) => b.className.includes('CloseButton')))

    expect(await screen.findByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should validate title', async () => {
    const user = userEvent.setup()

    reduxRender(<AddNewAlbumModal opened={true} onClose={() => {}} />)

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
