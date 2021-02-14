class MessagesController < ApplicationController
  before_action :set_chat
  before_action :set_message, only: [:show, :update]

  def index
    @messages = @chat.messages.all

    render json: @messages
  end

  def show
    render json: @message
  end

  def update
    if @message.update(content: message_params['message_content'])
      render json: @message
    else
      render json: @message.errors, status: :unprocessable_entity
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_chat
      app = App.find_by(token: params[:app_token])
      @chat = app.chats.find_by(number: params[:chat_number])
    end

    def set_message
      @message = @chat.messages.find_by(number: params[:number])
    end

    # Only allow a trusted parameter "white list" through.
    def message_params
      params.permit(:message_content)
    end
end
