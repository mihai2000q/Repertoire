import { reduxRender, withToastify } from '../../../test-utils.tsx'
import { setupServer } from 'msw/node'
import AddNewArtistModal from './AddNewArtistModal.tsx'
import { act, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { CreateArtistRequest } from '../../../types/requests/ArtistRequests.ts'

describe('Add New Artist Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    reduxRender(<AddNewArtistModal opened={true} onClose={() => {}} />)

    expect(screen.getByRole('dialog', { name: /add new artist/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /add new artist/i })).toBeInTheDocument()
    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /submit/i })).toBeInTheDocument()
  })

  it('should send only create request when no image is uploaded', async () => {
    const user = userEvent.setup()

    const newTitle = 'New Artist'

    const onClose = vitest.fn()

    let capturedRequest: CreateArtistRequest
    server.use(
      http.post('/artists', async (req) => {
        capturedRequest = (await req.request.json()) as CreateArtistRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(withToastify(<AddNewArtistModal opened={true} onClose={onClose} />))

    await user.type(screen.getByRole('textbox', { name: /name/i }), newTitle)
    await user.click(screen.getByRole('button', { name: /submit/i }))

    await waitFor(() =>
      expect(capturedRequest).toStrictEqual({
        name: newTitle
      })
    )
    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newTitle} added!`))
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
  })

  it('should send create request and save image request when the image is uploaded', async () => {
    const user = userEvent.setup()

    const newTitle = 'New Artist'
    const newImage = new File(['something'], 'image.png', { type: 'image/png' })

    const onClose = vitest.fn()

    const returnedId = 'the-artist-id'

    let capturedCreateRequest: CreateArtistRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.post('/artists', async (req) => {
        capturedCreateRequest = (await req.request.json()) as CreateArtistRequest
        return HttpResponse.json({ id: returnedId })
      }),
      http.put('/artists/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(withToastify(<AddNewArtistModal opened={true} onClose={onClose} />))

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    await user.type(screen.getByRole('textbox', { name: /name/i }), newTitle)
    await user.click(screen.getByRole('button', { name: /submit/i }))

    await waitFor(() =>
      expect(capturedCreateRequest).toStrictEqual({
        name: newTitle
      })
    )
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
    expect(capturedSaveImageFormData.get('id')).toBe(returnedId)

    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newTitle} added!`))
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
  })

  it('should remove the image on close', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })

    const onClose = vitest.fn()

    reduxRender(<AddNewArtistModal opened={true} onClose={onClose} />)

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()

    await user.click(screen.getAllByRole('button').find((b) => b.className.includes('CloseButton')))

    expect(await screen.findByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should validate name', async () => {
    const user = userEvent.setup()

    reduxRender(<AddNewArtistModal opened={true} onClose={() => {}} />)

    const name = screen.getByRole('textbox', { name: /name/i })

    expect(name).not.toBeInvalid()
    await user.click(screen.getByRole('button', { name: /submit/i }))
    expect(name).toBeInvalid()

    await user.type(name, 'something')
    expect(name).not.toBeInvalid()

    await user.clear(name)
    act(() => name.blur())
    expect(name).toBeInvalid()
  })
})
