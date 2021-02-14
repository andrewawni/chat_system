class ChatsController < ApplicationController
  before_action :set_app
  before_action :set_chat, only: [:show, :update]

  def index
    @chats = @app.chats.all

    render json: @chats
  end

  def show
    render json: @chat
  end

  def update
    if @chat.update(chat_params)
      render json: @chat
    else
      render json: @chat.errors, status: :unprocessable_entity
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_app
      @app = App.find_by(token: params[:app_token])
    end

    def set_chat
      @chat = @app.chats.find_by(number: params[:number])
    end

    # Only allow a trusted parameter "white list" through.
    def chat_params
      params.require(:chat).permit(:name)
    end
end
