package com.gusrylmubarok.springdi;

import com.gusrylmubarok.springdi.services.MessageService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.stereotype.Component;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Component
public class SetterBasedInjection {
    private MessageService messageService;

    @Autowired
    @Qualifier("TwitterService")
    public void setMessageService(MessageService messageService) {
        this.messageService = messageService;
    }

    public void processMsg(String message) {
        messageService.sendMsg(message);
    }
}
