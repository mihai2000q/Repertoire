import { reduxRender, withToastify } from '../../../../../../test-utils.tsx'
import AddNewSongArrangementButton from './AddNewSongArrangementButton.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event/dist/cjs/index.js'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { CreateSongArrangementRequest } from '../../../../../../types/requests/SongRequests.ts'

describe('Add New Song Arrangement Button', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    reduxRender(
      <AddNewSongArrangementButton openedPopover={false} setOpenedPopover={vi.fn()} songId={''} />
    )

    expect(screen.getByRole('button', { name: 'add-new-arrangement' })).toBeInTheDocument()
  })

  it('should render popover', async () => {
    reduxRender(
      <AddNewSongArrangementButton openedPopover={true} setOpenedPopover={vi.fn()} songId={''} />
    )

    expect(await screen.getByRole('dialog', { name: 'add-new-arrangement' })).toBeInTheDocument()
    expect(screen.getByText(/new arrangement/i)).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: 'name' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: 'name' })).not.toBeInvalid()
    expect(screen.getByRole('button', { name: 'submit' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'submit' })).not.toBeDisabled()
  })

  it('should send create request when the name has been entered and submitted', async () => {
    const user = userEvent.setup()

    const newName = 'Some Name'

    const songId = 'songId'
    const setOpenedPopover = vi.fn()
    const onCreate = vi.fn()

    const returnedId = 'the-arrangement-id'

    let capturedCreateRequest: CreateSongArrangementRequest
    server.use(
      http.post('/songs/arrangements', async (req) => {
        capturedCreateRequest = (await req.request.json()) as CreateSongArrangementRequest
        return HttpResponse.json({ id: returnedId })
      })
    )

    reduxRender(
      withToastify(
        <AddNewSongArrangementButton
          openedPopover={true}
          setOpenedPopover={setOpenedPopover}
          songId={songId}
          onCreate={onCreate}
        />
      )
    )

    await user.type(screen.getByRole('textbox', { name: 'name' }), newName)
    await user.click(screen.getByRole('button', { name: 'submit' }))

    expect(await screen.findByText(/new song arrangement added/i)).toBeInTheDocument()
    expect(capturedCreateRequest).toStrictEqual({ name: newName, songId: songId })
    expect(onCreate).toHaveBeenCalledExactlyOnceWith(returnedId)
    expect(setOpenedPopover).toHaveBeenCalledExactlyOnceWith(false)
    expect(screen.getByRole('textbox', { name: 'name' })).toHaveValue('')
  })

  it('should disable submit button when name field is invalid', async () => {
    const user = userEvent.setup()

    reduxRender(
      <AddNewSongArrangementButton openedPopover={true} setOpenedPopover={vi.fn()} songId={''} />
    )

    const nameTextbox = screen.getByRole('textbox', { name: 'name' })
    const submitButton = screen.getByRole('button', { name: 'submit' })

    await user.click(submitButton)
    expect(nameTextbox).toBeInvalid()
    expect(submitButton).toBeDisabled()

    await user.type(nameTextbox, 's')
    expect(nameTextbox).not.toBeInvalid()
    expect(submitButton).not.toBeDisabled()

    await user.clear(nameTextbox)
    expect(nameTextbox).toBeInvalid()
    expect(submitButton).toBeDisabled()
  })

  it('should validate name field', async () => {
    const user = userEvent.setup()

    reduxRender(
      <AddNewSongArrangementButton openedPopover={true} setOpenedPopover={vi.fn()} songId={''} />
    )

    const nameTextbox = screen.getByRole('textbox', { name: 'name' })

    await user.click(screen.getByRole('button', { name: 'submit' }))
    expect(nameTextbox).toBeInvalid()

    await user.type(nameTextbox, 's')
    expect(nameTextbox).not.toBeInvalid()

    await user.clear(nameTextbox)
    expect(nameTextbox).toBeInvalid()
  })
})
