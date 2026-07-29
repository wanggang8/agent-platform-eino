import { useState } from "react";
import type { FormEvent } from "react";
import { Send } from "lucide-react";

type ComposerProps = {
  readonly disabled: boolean;
  readonly onSubmit: (content: string) => Promise<void> | void;
};

// Composer 只提交用户文本，不在前端选择工具、状态或 fixture。
export function Composer({ disabled, onSubmit }: ComposerProps) {
  const [content, setContent] = useState("");
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const value = content.trim();
    if (!value || disabled) return;
    await onSubmit(value);
    setContent("");
  }
  return (
    <form className="chat-composer" data-testid="chat-composer" onSubmit={submit}>
      <label htmlFor="query-input">输入查询</label>
      <div>
        <textarea id="query-input" value={content} disabled={disabled} onChange={(event) => setContent(event.target.value)} placeholder="例如：查询新增的漏洞" rows={3} />
        <button type="submit" disabled={disabled || !content.trim()}><Send size={17} />发送</button>
      </div>
    </form>
  );
}
