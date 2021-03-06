class ChatsController < ApplicationController
  before_action :set_chat, only: [:show, :update]

  def index
    @chats = Chat.get_all(params[:app_token])

    render json: @chats
  end

  def show
    render json: @chat
  end

  def update
    if @chat.update(name: chat_params['chat_name'])
      render json: @chat
    else
      render json: @chat.errors, status: :unprocessable_entity
    end
  end

  def search
    chat = @app.chats.find_by(number: params[:chat_number])
    q = if params['search_query']
          params['search_query'] + '*'
        else
          '*'
        end
    @messages_results = chat.messages.search(q).records if chat
    render json: @messages_results
  end

  private

  # Use callbacks to share common setup or constraints between actions.
  def set_chat
    @chat = Chat.get_one(params[:app_token], params[:number])
  end

  # Only allow a trusted parameter "white list" through.
  def chat_params
    params.permit('chat_name')
  end
end
