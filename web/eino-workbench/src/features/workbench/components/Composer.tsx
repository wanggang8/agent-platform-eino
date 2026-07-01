import { Paperclip, Send, Wrench } from "lucide-react";

type ComposerProps = {
  readonly disabled: boolean;
};

// Composer 是消息输入区；disabled 时只表示等待审批/澄清，不在前端推进 run 状态。
export function Composer({ disabled }: ComposerProps) {
  return (
    <form className="chat-composer" data-testid="chat-composer">
      <div className="mode-row" aria-label="模式">
        <button type="button" className="is-active">智能模式</button>
        <button type="button">分析模式</button>
        <button type="button">执行模式</button>
        <button type="button">只读模式</button>
      </div>
      <label>
        <span className="sr-only">输入你的问题</span>
        <textarea disabled={disabled} placeholder={disabled ? "等待审批或澄清后继续" : "输入你的问题，或使用 / 选择能力"} />
      </label>
      <div className="composer-actions">
        <button type="button"><Wrench size={15} /> 工具</button>
        <button type="button">能力</button>
        <button type="button"><Paperclip size={15} /> 附件</button>
        <button type="submit" className="send-button" disabled={disabled} aria-label="发送"><Send size={16} /></button>
      </div>
    </form>
  );
}
