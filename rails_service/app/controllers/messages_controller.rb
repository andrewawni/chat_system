# frozen_string_literal: true

class MessagesController < ApplicationController
  prepend_before_action :set_message
  before_action :set_message, only: [:show, :update]

  def index
    @messages = Message.get_all(params[:app_token], params[:chat_number])
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
  def set_message
   @message = Message.get_one(params[:app_token], params[:chat_number], params[:number])
  end

  # Only allow a trusted parameter "white list" through.
  def message_params
    params.permit(:message_content)
  end
end
